package postgres

import (
	"database/sql"
	"time"
)

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

func (r *TransactionRepository) List() ([]Transaction, error) {
	rows, err := r.db.Query("SELECT id, order_id, amount, payment_method, status, created_at FROM transactions ORDER BY id")
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
	return list, nil
}

func (r *TransactionRepository) GetByID(id int64) (*Transaction, error) {
	t := &Transaction{}
	err := r.db.QueryRow("SELECT id, order_id, amount, payment_method, status, created_at FROM transactions WHERE id=$1", id).Scan(&t.ID, &t.OrderID, &t.Amount, &t.PaymentMethod, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *TransactionRepository) Create(orderID int64, amount float64, paymentMethod string) (*Transaction, error) {
	t := &Transaction{}
	err := r.db.QueryRow("INSERT INTO transactions (order_id, amount, payment_method) VALUES ($1,$2,$3) RETURNING id, order_id, amount, payment_method, status, created_at", orderID, amount, paymentMethod).Scan(&t.ID, &t.OrderID, &t.Amount, &t.PaymentMethod, &t.Status, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}
