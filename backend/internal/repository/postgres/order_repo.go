package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

type Order struct {
	ID             int64     `json:"id"`
	OrderNumber    string    `json:"order_number"`
	ShiftID        int64     `json:"shift_id"`
	UserID         int64     `json:"user_id"`
	OrderType      string    `json:"order_type"`
	TableNumber    string    `json:"table_number"`
	CustomerName   string    `json:"customer_name"`
	CustomerPhone  string    `json:"customer_phone"`
	Total          float64   `json:"total"`
	DiscountAmount float64   `json:"discount_amount"`
	PaymentMethod  string    `json:"payment_method"`
	OutletID       string    `json:"outlet_id"`
	Status         string    `json:"status"`
	ClientOrderID  string    `json:"client_order_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateOrderInput struct {
	ShiftID             *int64
	UserID              int64
	OrderType           string
	TableNumber         string
	CustomerName        string
	CustomerPhone       string
	Items               []OrderItemInput
	Total               float64
	DiscountAmount      float64
	DiscountPercent     float64
	DiscountCode        string
	PaymentMethod       string
	OutletID            string
	ClientOrderID       string
	Payments            []PaymentInput
	SkipShiftValidation bool
}

type PaymentInput struct {
	Amount float64
	Method string
	RefNo  string
}

type OrderItemInput struct {
	ProductID    int64
	Quantity     int
	Price        float64
	Subtotal     float64
	VariantLabel string
	Note         string
}

type OrderResult struct {
	Order         *Order
	OrderItems    []OrderItem
	AlreadySynced bool
	Items         []OrderItem
	Payments      []PaymentInput
	Total         float64
	Subtotal      float64
	Discount      float64
	Tax           float64
	Change        float64
	CreatedAt     time.Time
	StoreName     string
	StoreAddr     string
	StorePhone    string
	CashierName   string
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

const orderCols = "id, order_number, shift_id, user_id, order_type, table_number, customer_name, customer_phone, total, discount_amount, payment_method, outlet_id, status, client_order_id, created_at, updated_at"

func scanOrder(o *Order, row scanRow) error {
	return row.Scan(&o.ID, &o.OrderNumber, &o.ShiftID, &o.UserID, &o.OrderType, &o.TableNumber, &o.CustomerName, &o.CustomerPhone, &o.Total, &o.DiscountAmount, &o.PaymentMethod, &o.OutletID, &o.Status, &o.ClientOrderID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*Order, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Order, error) {
		o := &Order{}
		err := tx.QueryRowContext(ctx, "SELECT "+orderCols+" FROM orders WHERE id=$1 AND tenant_id=$2", id, tenantID).Scan(
			&o.ID, &o.OrderNumber, &o.ShiftID, &o.UserID, &o.OrderType, &o.TableNumber,
			&o.CustomerName, &o.CustomerPhone, &o.Total, &o.DiscountAmount, &o.PaymentMethod,
			&o.OutletID, &o.Status, &o.ClientOrderID, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return o, nil
	})
}

// posTaxRate mirrors the cashier UI: tax = (subtotal - discount) * 11%.
const posTaxRate = 0.11

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// CreateOrder persists an order atomically: idempotency check (client_order_id),
// catalog pricing, stock validation, orders + order_items + transactions rows,
// and stock deduction. Returns the already-persisted order when the same
// client_order_id arrives again (offline sync retry).
func (r *OrderRepository) CreateOrder(ctx context.Context, in CreateOrderInput) (*OrderResult, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*OrderResult, error) {
		now := time.Now()

		if in.ClientOrderID != "" {
			var existingID int64
			err := tx.QueryRowContext(ctx,
				"SELECT id FROM orders WHERE client_order_id=$1 AND tenant_id=$2",
				in.ClientOrderID, tenantID).Scan(&existingID)
			if err == nil {
				res, lerr := r.loadOrderResultTx(ctx, tx, tenantID, existingID)
				if lerr != nil {
					return nil, lerr
				}
				res.AlreadySynced = true
				return res, nil
			}
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		if len(in.Items) == 0 {
			return nil, errors.New("order tidak memiliki item")
		}
		if in.OrderType == "" {
			in.OrderType = "dine_in"
		}

		type orderLine struct {
			item      OrderItemInput
			name      string
			price     float64
			subtotal  float64
			isService bool
		}
		type productInfo struct {
			name      string
			basePrice float64
			isService bool
			available int
		}

		products := map[int64]*productInfo{}
		requiredQty := map[int64]int{}
		var lines []orderLine
		var gross float64

		for _, it := range in.Items {
			if it.ProductID <= 0 {
				return nil, fmt.Errorf("%w: product_id tidak valid", ErrProductNotFound)
			}
			if it.Quantity <= 0 {
				return nil, errors.New("jumlah item harus lebih dari 0")
			}
			info, ok := products[it.ProductID]
			if !ok {
				info = &productInfo{}
				var stockCol int
				var available float64
				err := tx.QueryRowContext(ctx,
					`SELECT p.name, p.base_price, p.is_service, p.stock,
						COALESCE((SELECT SUM(s.quantity) FROM stock s WHERE s.product_id=p.id AND s.tenant_id=p.tenant_id), p.stock)
					 FROM products p
					 WHERE p.id=$1 AND p.tenant_id=$2`,
					it.ProductID, tenantID).Scan(&info.name, &info.basePrice, &info.isService, &stockCol, &available)
				if err == sql.ErrNoRows {
					return nil, fmt.Errorf("%w: produk #%d", ErrProductNotFound, it.ProductID)
				}
				if err != nil {
					return nil, err
				}
				info.available = int(available)
				products[it.ProductID] = info
			}

			price := it.Price
			if price <= 0 {
				price = info.basePrice
			}
			sub := it.Subtotal
			if sub <= 0 {
				sub = price * float64(it.Quantity)
			}
			gross += sub
			requiredQty[it.ProductID] += it.Quantity
			lines = append(lines, orderLine{
				item:      it,
				name:      info.name,
				price:     price,
				subtotal:  sub,
				isService: info.isService,
			})
		}

		for pid, need := range requiredQty {
			info := products[pid]
			if info.isService {
				continue
			}
			if info.available < need {
				return nil, fmt.Errorf("%w: %s", ErrInsufficientStock, info.name)
			}
		}

		discount := in.DiscountAmount
		if discount <= 0 && in.DiscountPercent > 0 {
			discount = gross * in.DiscountPercent / 100
		}
		if discount < 0 {
			discount = 0
		}
		if discount > gross {
			discount = gross
		}
		tax := (gross - discount) * posTaxRate
		total := gross - discount + tax
		gross, discount, tax, total = round2(gross), round2(discount), round2(tax), round2(total)

		paymentMethod := in.PaymentMethod
		if paymentMethod == "" && len(in.Payments) > 0 {
			paymentMethod = in.Payments[0].Method
		}
		if paymentMethod == "" {
			paymentMethod = "cash"
		}

		if in.ShiftID != nil && !in.SkipShiftValidation {
			var status string
			err := tx.QueryRowContext(ctx,
				"SELECT status FROM shifts WHERE id=$1 AND tenant_id=$2",
				*in.ShiftID, tenantID).Scan(&status)
			if err == sql.ErrNoRows {
				return nil, ErrInvalidShift
			}
			if err != nil {
				return nil, err
			}
			if status != "open" {
				return nil, ErrShiftNotOpen
			}
		}

		var orderID int64
		err := tx.QueryRowContext(ctx,
			`INSERT INTO orders (tenant_id, order_number, shift_id, user_id, order_type, table_number,
				customer_name, customer_phone, total, discount_amount, payment_method, outlet_id,
				status, client_order_id, created_at, updated_at)
			 VALUES ($1,'',$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,'completed',$12,$13,$13)
			 RETURNING id`,
			tenantID, in.ShiftID, in.UserID, in.OrderType, in.TableNumber,
			in.CustomerName, in.CustomerPhone, total, discount, paymentMethod, in.OutletID,
			in.ClientOrderID, now).Scan(&orderID)
		if err != nil {
			return nil, err
		}

		orderNumber := fmt.Sprintf("ORD-%06d", orderID)
		if _, err := tx.ExecContext(ctx,
			"UPDATE orders SET order_number=$1 WHERE id=$2 AND tenant_id=$3",
			orderNumber, orderID, tenantID); err != nil {
			return nil, err
		}

		resultItems := make([]OrderItem, 0, len(lines))
		for _, ln := range lines {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO order_items (order_id, product_id, quantity, price, subtotal, variant_label, note, tenant_id)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
				orderID, ln.item.ProductID, ln.item.Quantity, ln.price, ln.subtotal,
				ln.item.VariantLabel, ln.item.Note, tenantID); err != nil {
				return nil, err
			}
			resultItems = append(resultItems, OrderItem{
				OrderID:      orderID,
				ProductID:    ln.item.ProductID,
				ProductName:  ln.name,
				VariantLabel: ln.item.VariantLabel,
				Note:         ln.item.Note,
				Quantity:     ln.item.Quantity,
				Price:        ln.price,
				Subtotal:     ln.subtotal,
			})
		}

		resultPayments := in.Payments
		if len(resultPayments) == 0 {
			resultPayments = []PaymentInput{{Method: paymentMethod, Amount: total}}
		}
		var paidTotal float64
		for _, p := range resultPayments {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO transactions (order_id, amount, payment_method, tenant_id) VALUES ($1,$2,$3,$4)",
				orderID, p.Amount, p.Method, tenantID); err != nil {
				return nil, err
			}
			paidTotal += p.Amount
		}

		var warehouseID int64
		whErr := tx.QueryRowContext(ctx,
			"SELECT id FROM warehouses WHERE is_active=true AND tenant_id=$1 ORDER BY id LIMIT 1",
			tenantID).Scan(&warehouseID)
		if whErr != nil && whErr != sql.ErrNoRows {
			return nil, whErr
		}
		hasWarehouse := whErr == nil
		for pid, need := range requiredQty {
			info := products[pid]
			if info.isService || need == 0 {
				continue
			}
			if hasWarehouse {
				if _, err := tx.ExecContext(ctx,
					`INSERT INTO stock (product_id, warehouse_id, quantity, tenant_id)
					 VALUES ($1,$2,$3,$4)
					 ON CONFLICT (product_id, warehouse_id)
					 DO UPDATE SET quantity = stock.quantity + $3, updated_at = NOW()`,
					pid, warehouseID, -need, tenantID); err != nil {
					return nil, err
				}
				if _, err := tx.ExecContext(ctx,
					`INSERT INTO stock_movements (product_id, warehouse_id, quantity, movement_type, reference_type, reference_id, note, created_by, tenant_id)
					 VALUES ($1,$2,$3,'sale','order',$4,$5,$6,$7)`,
					pid, warehouseID, -need, orderID, "penjualan POS", in.UserID, tenantID); err != nil {
					return nil, err
				}
				if _, err := tx.ExecContext(ctx,
					`UPDATE products SET stock = COALESCE((SELECT SUM(quantity) FROM stock WHERE product_id=$1 AND tenant_id=$2), 0)
					 WHERE id=$1 AND tenant_id=$2`,
					pid, tenantID); err != nil {
					return nil, err
				}
			} else {
				if _, err := tx.ExecContext(ctx,
					"UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id=$2 AND tenant_id=$3",
					need, pid, tenantID); err != nil {
					return nil, err
				}
			}
		}

		var storeName, cashierName string
		_ = tx.QueryRowContext(ctx, "SELECT name FROM tenants WHERE id=$1", tenantID).Scan(&storeName)
		_ = tx.QueryRowContext(ctx,
			"SELECT name FROM users WHERE id=$1 AND tenant_id=$2::text",
			in.UserID, tenantID).Scan(&cashierName)

		order := &Order{
			ID:             orderID,
			OrderNumber:    orderNumber,
			UserID:         in.UserID,
			OrderType:      in.OrderType,
			TableNumber:    in.TableNumber,
			CustomerName:   in.CustomerName,
			CustomerPhone:  in.CustomerPhone,
			Total:          total,
			DiscountAmount: discount,
			PaymentMethod:  paymentMethod,
			OutletID:       in.OutletID,
			Status:         "completed",
			ClientOrderID:  in.ClientOrderID,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if in.ShiftID != nil {
			order.ShiftID = *in.ShiftID
		}

		change := round2(paidTotal - total)
		if change < 0 {
			change = 0
		}

		return &OrderResult{
			Order:       order,
			Items:       resultItems,
			Payments:    resultPayments,
			Total:       total,
			Subtotal:    gross,
			Discount:    discount,
			Tax:         tax,
			Change:      change,
			CreatedAt:   now,
			StoreName:   storeName,
			CashierName: cashierName,
		}, nil
	})
}

// loadOrderResultTx rebuilds a full OrderResult from persisted rows
// (used when an offline sync retries an already-saved client_order_id).
func (r *OrderRepository) loadOrderResultTx(ctx context.Context, tx *sql.Tx, tenantID, orderID int64) (*OrderResult, error) {
	order := &Order{}
	err := tx.QueryRowContext(ctx,
		`SELECT id, COALESCE(order_number,''), COALESCE(shift_id,0), user_id, COALESCE(order_type,''),
			COALESCE(table_number,''), COALESCE(customer_name,''), COALESCE(customer_phone,''),
			total, COALESCE(discount_amount,0), COALESCE(payment_method,''), COALESCE(outlet_id,''),
			COALESCE(status,''), COALESCE(client_order_id,''), created_at, updated_at
		 FROM orders WHERE id=$1 AND tenant_id=$2`, orderID, tenantID).Scan(
		&order.ID, &order.OrderNumber, &order.ShiftID, &order.UserID, &order.OrderType,
		&order.TableNumber, &order.CustomerName, &order.CustomerPhone, &order.Total,
		&order.DiscountAmount, &order.PaymentMethod, &order.OutletID, &order.Status,
		&order.ClientOrderID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT oi.id, oi.order_id, oi.product_id, COALESCE(p.name,''), COALESCE(oi.variant_label,''),
			COALESCE(oi.note,''), oi.quantity, oi.price, oi.subtotal
		 FROM order_items oi
		 LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id = oi.tenant_id
		 WHERE oi.order_id=$1 AND oi.tenant_id=$2
		 ORDER BY oi.id`, orderID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]OrderItem, 0)
	var gross float64
	for rows.Next() {
		var oi OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.ProductName, &oi.VariantLabel,
			&oi.Note, &oi.Quantity, &oi.Price, &oi.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, oi)
		gross += oi.Subtotal
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	payRows, err := tx.QueryContext(ctx,
		`SELECT amount, COALESCE(payment_method,'') FROM transactions
		 WHERE order_id=$1 AND tenant_id=$2 ORDER BY id`, orderID, tenantID)
	if err != nil {
		return nil, err
	}
	defer payRows.Close()
	payments := make([]PaymentInput, 0)
	var paidTotal float64
	for payRows.Next() {
		var p PaymentInput
		if err := payRows.Scan(&p.Amount, &p.Method); err != nil {
			return nil, err
		}
		payments = append(payments, p)
		paidTotal += p.Amount
	}
	if err := payRows.Err(); err != nil {
		return nil, err
	}

	var storeName, cashierName string
	_ = tx.QueryRowContext(ctx, "SELECT name FROM tenants WHERE id=$1", tenantID).Scan(&storeName)
	_ = tx.QueryRowContext(ctx,
		"SELECT name FROM users WHERE id=$1 AND tenant_id=$2::text",
		order.UserID, tenantID).Scan(&cashierName)

	discount := order.DiscountAmount
	tax := round2(order.Total - round2(gross) + discount)
	change := round2(paidTotal - order.Total)
	if change < 0 {
		change = 0
	}

	return &OrderResult{
		Order:       order,
		Items:       items,
		Payments:    payments,
		Total:       order.Total,
		Subtotal:    round2(gross),
		Discount:    discount,
		Tax:         tax,
		Change:      change,
		CreatedAt:   order.CreatedAt,
		StoreName:   storeName,
		CashierName: cashierName,
	}, nil
}

func (r *OrderRepository) RecentOrders(ctx context.Context, limit int) ([]RecentOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]RecentOrder, error) {
		if limit <= 0 {
			limit = 10
		}
		rows, err := tx.QueryContext(ctx, `
			SELECT o.id, o.order_number, o.total, COALESCE(o.payment_method,''), COALESCE(u.name,''),
				COALESCE((SELECT SUM(oi.quantity) FROM order_items oi WHERE oi.order_id=o.id AND oi.tenant_id=o.tenant_id),0),
				o.status, TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI:SS')
			FROM orders o
			LEFT JOIN users u ON u.id = o.user_id AND u.tenant_id = o.tenant_id::text
			WHERE o.tenant_id=$1
			ORDER BY o.created_at DESC
			LIMIT $2`, tenantID, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]RecentOrder, 0)
		for rows.Next() {
			var ro RecentOrder
			if err := rows.Scan(&ro.ID, &ro.OrderNumber, &ro.Total, &ro.PaymentMethod, &ro.CashierName, &ro.ItemCount, &ro.Status, &ro.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, ro)
		}
		return list, rows.Err()
	})
}

func (r *OrderRepository) RevenueSummary(ctx context.Context, period string) (map[string]any, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (map[string]any, error) {
		return map[string]any{}, nil
	})
}

func (r *OrderRepository) GetDefaultWarehouseID(ctx context.Context) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		return 0, nil
	})
}

func (r *OrderRepository) OpenShift(ctx context.Context, userID int64, openingBalance float64) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO shifts (tenant_id, user_id, opening_balance, status) VALUES ($1,$2,$3,'open') RETURNING id",
			tenantID, userID, openingBalance).Scan(&id)
		return id, err
	})
}

func (r *OrderRepository) CloseShift(ctx context.Context, shiftID int64, closingBalance float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE shifts SET closing_balance=$1, closed_at=NOW(), status='closed' WHERE id=$2 AND tenant_id=$3",
			closingBalance, shiftID, tenantID)
		return err
	})
}

var (
	ErrDuplicateOrder    = errors.New("duplicate order")
	ErrInvalidShift      = errors.New("invalid shift")
	ErrShiftNotOpen      = errors.New("shift not open")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)
