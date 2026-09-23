package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Order struct {
	ID              int64      `json:"id"`
	OrderNumber     string     `json:"order_number"`
	ShiftID         int64      `json:"shift_id"`
	UserID          int64      `json:"user_id"`
	OrderType       string     `json:"order_type"`
	TableNumber     string     `json:"table_number"`
	CustomerName    string     `json:"customer_name"`
	CustomerPhone   string     `json:"customer_phone"`
	Total           float64    `json:"total"`
	DiscountAmount  float64    `json:"discount_amount"`
	PaymentMethod   string     `json:"payment_method"`
	OutletID        string     `json:"outlet_id"`
	Status          string     `json:"status"`
	ClientOrderID   string     `json:"client_order_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
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
	DiscountCode        string
	PaymentMethod       string
	OutletID            string
	ClientOrderID       string
	Payments            []PaymentInput
	SkipShiftValidation bool
}

type PaymentInput struct {
	Amount    float64
	Method    string
	RefNo     string
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
	StoreName     string
	StoreAddr     string
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
		err := tx.QueryRowContext(ctx, "SELECT "+orderCols+" FROM orders WHERE id=$1", id).Scan(
			&o.ID, &o.OrderNumber, &o.ShiftID, &o.UserID, &o.OrderType, &o.TableNumber,
			&o.CustomerName, &o.CustomerPhone, &o.Total, &o.DiscountAmount, &o.PaymentMethod,
			&o.OutletID, &o.Status, &o.ClientOrderID, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return o, nil
	})
}

func (r *OrderRepository) CreateOrder(ctx context.Context, in CreateOrderInput) (*OrderResult, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*OrderResult, error) {
		// Simplified implementation - just return minimal result
		return &OrderResult{Order: &Order{ID: 1}}, nil
	})
}

func (r *OrderRepository) RecentOrders(ctx context.Context, limit int) ([]RecentOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]RecentOrder, error) {
		return []RecentOrder{}, nil
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
		return 0, nil
	})
}

func (r *OrderRepository) CloseShift(ctx context.Context, shiftID int64, closingBalance float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		return nil
	})
}

var (
	ErrDuplicateOrder    = errors.New("duplicate order")
	ErrInvalidShift      = errors.New("invalid shift")
	ErrShiftNotOpen      = errors.New("shift not open")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)