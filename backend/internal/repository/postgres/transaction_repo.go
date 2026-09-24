package postgres

import (
	"context"
	"database/sql"
	"time"
)

type TransactionLog struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	OrderID       int64     `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionLogItem struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

type Transaction struct {
	ID            int64     `json:"id"`
	OrderID       int64     `json:"order_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) List(ctx context.Context) ([]Transaction, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Transaction, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, order_id, amount, payment_method, status, created_at FROM transactions WHERE tenant_id=$1 ORDER BY id", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var list []Transaction
		for rows.Next() {
			var t Transaction
			if err := rows.Scan(&t.ID, &t.OrderID, &t.Amount, &t.PaymentMethod, &t.Status, &t.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, t)
		}
		return list, rows.Err()
	})
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Transaction, error) {
		t := &Transaction{}
		err := tx.QueryRowContext(ctx, "SELECT id, order_id, amount, payment_method, status, created_at FROM transactions WHERE id=$1 AND tenant_id=$2", id, tenantID).Scan(&t.ID, &t.OrderID, &t.Amount, &t.PaymentMethod, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		return t, nil
	})
}

func (r *TransactionRepository) Create(ctx context.Context, orderID int64, amount float64, paymentMethod string) (*Transaction, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Transaction, error) {
		t := &Transaction{}
		err := tx.QueryRowContext(ctx, "INSERT INTO transactions (order_id, amount, payment_method, tenant_id) VALUES ($1,$2,$3,$4) RETURNING id, order_id, amount, payment_method, status, created_at", orderID, amount, paymentMethod, tenantID).Scan(&t.ID, &t.OrderID, &t.Amount, &t.PaymentMethod, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		return t, nil
	})
}

func (r *TransactionRepository) GetLogsByUsername(ctx context.Context, username string) ([]TransactionLog, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]TransactionLog, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT t.id, COALESCE(u.username,''), t.order_id, t.amount, t.payment_method, t.status, t.created_at
			FROM transactions t
			LEFT JOIN orders o ON o.id = t.order_id AND o.tenant_id = t.tenant_id
			LEFT JOIN users u ON u.id = o.user_id AND u.tenant_id = t.tenant_id::text
			WHERE t.tenant_id=$1 AND u.username=$2
			ORDER BY t.created_at DESC`, tenantID, username)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]TransactionLog, 0)
		for rows.Next() {
			var l TransactionLog
			if err := rows.Scan(&l.ID, &l.Username, &l.OrderID, &l.Amount, &l.PaymentMethod, &l.Status, &l.CreatedAt); err != nil {
				return nil, err
			}
			list = append(list, l)
		}
		return list, rows.Err()
	})
}

func (r *TransactionRepository) GetLog(ctx context.Context, id int64) (*TransactionLog, error) {
	return r.GetLogByID(ctx, id)
}

func (r *TransactionRepository) GetLogByID(ctx context.Context, id int64) (*TransactionLog, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*TransactionLog, error) {
		l := &TransactionLog{}
		err := tx.QueryRowContext(ctx, `
			SELECT t.id, COALESCE(u.username,''), t.order_id, t.amount, t.payment_method, t.status, t.created_at
			FROM transactions t
			LEFT JOIN orders o ON o.id = t.order_id AND o.tenant_id = t.tenant_id
			LEFT JOIN users u ON u.id = o.user_id AND u.tenant_id = t.tenant_id::text
			WHERE t.id=$1 AND t.tenant_id=$2`, id, tenantID).
			Scan(&l.ID, &l.Username, &l.OrderID, &l.Amount, &l.PaymentMethod, &l.Status, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		return l, nil
	})
}
