package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type MarketplaceConnection struct {
	ID                 int64     `json:"id"`
	Platform           string    `json:"platform"`
	ShopName           string    `json:"shop_name"`
	APIToken           string    `json:"api_token"`
	AccountName        string    `json:"account_name"`
	CustomerID         *int64    `json:"customer_id"`
	CustomerName       string    `json:"customer_name"`
	SalesCategoryID    *int64    `json:"sales_category_id"`
	SalesCategoryName  string    `json:"sales_category_name"`
	ShippingFeeAccount string    `json:"shipping_fee_account"`
	CommissionAccount  string    `json:"commission_account"`
	IsActive           bool      `json:"is_active"`
	LastSyncAt         *string   `json:"last_sync_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type MarketplaceOrder struct {
	ID                 int64     `json:"id"`
	ConnectionID       int64     `json:"connection_id"`
	Platform           string    `json:"platform"`
	ShopName           string    `json:"shop_name"`
	MarketplaceOrderID string    `json:"marketplace_order_id"`
	OrderDate          string    `json:"order_date"`
	CustomerName       string    `json:"customer_name"`
	ProductID          *int64    `json:"product_id"`
	ProductName        string    `json:"product_name"`
	Quantity           float64   `json:"quantity"`
	UnitPrice          float64   `json:"unit_price"`
	Subtotal           float64   `json:"subtotal"`
	ShippingFee        float64   `json:"shipping_fee"`
	PlatformFee        float64   `json:"platform_fee"`
	GrandTotal         float64   `json:"grand_total"`
	Status             string    `json:"status"`
	SalesOrderID       *int64    `json:"sales_order_id"`
	SalesOrderNumber   string    `json:"sales_order_number"`
	CreatedBy          *int64    `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
}

type MarketplaceRepository struct {
	db *sql.DB
}

func NewMarketplaceRepository(db *sql.DB) *MarketplaceRepository {
	return &MarketplaceRepository{db: db}
}

func (r *MarketplaceRepository) ListConnections(ctx context.Context) ([]MarketplaceConnection, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]MarketplaceConnection, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT c.id, c.platform, c.shop_name, c.api_token, c.account_name,
				c.customer_id, COALESCE(cu.name, ''), c.sales_category_id, COALESCE(sc.name, ''),
				c.shipping_fee_account, c.commission_account, c.is_active,
				c.last_sync_at::text, c.created_at, c.updated_at
			FROM marketplace_connections c
			LEFT JOIN customers cu ON cu.id = c.customer_id
			LEFT JOIN sales_categories sc ON sc.id = c.sales_category_id
			WHERE c.tenant_id=$1
			ORDER BY c.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]MarketplaceConnection, 0)
		for rows.Next() {
			var c MarketplaceConnection
			if err := rows.Scan(&c.ID, &c.Platform, &c.ShopName, &c.APIToken, &c.AccountName,
				&c.CustomerID, &c.CustomerName, &c.SalesCategoryID, &c.SalesCategoryName,
				&c.ShippingFeeAccount, &c.CommissionAccount, &c.IsActive,
				&c.LastSyncAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, c)
		}
		return list, nil
	})
}

func (r *MarketplaceRepository) GetConnection(ctx context.Context, id int64) (*MarketplaceConnection, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*MarketplaceConnection, error) {
		var c MarketplaceConnection
		err := tx.QueryRowContext(ctx, `
			SELECT c.id, c.platform, c.shop_name, c.api_token, c.account_name,
				c.customer_id, COALESCE(cu.name, ''), c.sales_category_id, COALESCE(sc.name, ''),
				c.shipping_fee_account, c.commission_account, c.is_active,
				c.last_sync_at::text, c.created_at, c.updated_at
			FROM marketplace_connections c
			LEFT JOIN customers cu ON cu.id = c.customer_id
			LEFT JOIN sales_categories sc ON sc.id = c.sales_category_id
			WHERE c.id=$1 AND c.tenant_id=$2`, id, tenantID).
			Scan(&c.ID, &c.Platform, &c.ShopName, &c.APIToken, &c.AccountName,
				&c.CustomerID, &c.CustomerName, &c.SalesCategoryID, &c.SalesCategoryName,
				&c.ShippingFeeAccount, &c.CommissionAccount, &c.IsActive,
				&c.LastSyncAt, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &c, nil
	})
}

func (r *MarketplaceRepository) CreateConnection(ctx context.Context, c *MarketplaceConnection) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		err := tx.QueryRowContext(ctx, `
			INSERT INTO marketplace_connections (platform, shop_name, api_token, account_name,
				customer_id, sales_category_id, shipping_fee_account, commission_account, is_active, tenant_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
			c.Platform, c.ShopName, c.APIToken, c.AccountName,
			c.CustomerID, c.SalesCategoryID, c.ShippingFeeAccount, c.CommissionAccount, c.IsActive, tenantID).Scan(&c.ID)
		if err != nil {
			return 0, err
		}
		return c.ID, nil
	})
}

func (r *MarketplaceRepository) UpdateConnection(ctx context.Context, c *MarketplaceConnection) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, `
			UPDATE marketplace_connections SET platform=$1, shop_name=$2, api_token=$3, account_name=$4,
				customer_id=$5, sales_category_id=$6, shipping_fee_account=$7, commission_account=$8,
				is_active=$9, updated_at=NOW() WHERE id=$10 AND tenant_id=$11`,
			c.Platform, c.ShopName, c.APIToken, c.AccountName,
			c.CustomerID, c.SalesCategoryID, c.ShippingFeeAccount, c.CommissionAccount,
			c.IsActive, c.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) DeleteConnection(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM marketplace_connections WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) UpdateLastSync(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE marketplace_connections SET last_sync_at=NOW() WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) ListOrders(ctx context.Context) ([]MarketplaceOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]MarketplaceOrder, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT o.id, o.connection_id, c.platform, c.shop_name, o.marketplace_order_id,
				o.order_date::text, o.customer_name, o.product_id, o.product_name, o.quantity,
				o.unit_price, o.subtotal, o.shipping_fee, o.platform_fee, o.grand_total,
				o.status, o.sales_order_id, COALESCE(so.order_number, ''), o.created_by, o.created_at
			FROM marketplace_orders o
			JOIN marketplace_connections c ON c.id = o.connection_id
			LEFT JOIN sales_orders so ON so.id = o.sales_order_id
			WHERE o.tenant_id=$1
			ORDER BY o.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]MarketplaceOrder, 0)
		for rows.Next() {
			var o MarketplaceOrder
			if err := rows.Scan(&o.ID, &o.ConnectionID, &o.Platform, &o.ShopName, &o.MarketplaceOrderID,
				&o.OrderDate, &o.CustomerName, &o.ProductID, &o.ProductName, &o.Quantity,
				&o.UnitPrice, &o.Subtotal, &o.ShippingFee, &o.PlatformFee, &o.GrandTotal,
				&o.Status, &o.SalesOrderID, &o.SalesOrderNumber, &o.CreatedBy, &o.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, o)
		}
		return list, nil
	})
}

func (r *MarketplaceRepository) GetOrder(ctx context.Context, id int64) (*MarketplaceOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*MarketplaceOrder, error) {
		var o MarketplaceOrder
		err := tx.QueryRowContext(ctx, `
			SELECT o.id, o.connection_id, c.platform, c.shop_name, o.marketplace_order_id,
				o.order_date::text, o.customer_name, o.product_id, o.product_name, o.quantity,
				o.unit_price, o.subtotal, o.shipping_fee, o.platform_fee, o.grand_total,
				o.status, o.sales_order_id, COALESCE(so.order_number, ''), o.created_by, o.created_at
			FROM marketplace_orders o
			JOIN marketplace_connections c ON c.id = o.connection_id
			LEFT JOIN sales_orders so ON so.id = o.sales_order_id
			WHERE o.id=$1 AND o.tenant_id=$2`, id, tenantID).
			Scan(&o.ID, &o.ConnectionID, &o.Platform, &o.ShopName, &o.MarketplaceOrderID,
				&o.OrderDate, &o.CustomerName, &o.ProductID, &o.ProductName, &o.Quantity,
				&o.UnitPrice, &o.Subtotal, &o.ShippingFee, &o.PlatformFee, &o.GrandTotal,
				&o.Status, &o.SalesOrderID, &o.SalesOrderNumber, &o.CreatedBy, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &o, nil
	})
}

func (r *MarketplaceRepository) OrderExists(ctx context.Context, connectionID int64, marketplaceOrderID string) (bool, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (bool, error) {
		var exists bool
		err := tx.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM marketplace_orders WHERE connection_id=$1 AND marketplace_order_id=$2 AND tenant_id=$3)",
			connectionID, marketplaceOrderID, tenantID).Scan(&exists)
		if err != nil {
			return false, err
		}
		return exists, nil
	})
}

func (r *MarketplaceRepository) CreateOrder(ctx context.Context, o *MarketplaceOrder) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO marketplace_orders (connection_id, marketplace_order_id, order_date, customer_name,
				product_id, product_name, quantity, unit_price, subtotal, shipping_fee, platform_fee,
				grand_total, status, created_by, tenant_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
			o.ConnectionID, o.MarketplaceOrderID, o.OrderDate, o.CustomerName,
			o.ProductID, o.ProductName, o.Quantity, o.UnitPrice, o.Subtotal, o.ShippingFee, o.PlatformFee,
			o.GrandTotal, o.Status, o.CreatedBy, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) LinkSalesOrder(ctx context.Context, orderID int64, salesOrderID int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx,
			"UPDATE marketplace_orders SET status='posted', sales_order_id=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3",
			salesOrderID, orderID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) DeleteOrder(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM marketplace_orders WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *MarketplaceRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM marketplace_orders WHERE tenant_id=$1", tenantID).Scan(&count); err != nil {
			return "", err
		}
		return fmt.Sprintf("MP-%05d", count+1), nil
	})
}
