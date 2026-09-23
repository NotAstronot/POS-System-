package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PriceChangeItem struct {
	ID            int64   `json:"id"`
	PriceChangeID int64   `json:"price_change_id"`
	ProductID     int64   `json:"product_id"`
	ProductName   string  `json:"product_name"`
	OldPrice      float64 `json:"old_price"`
	NewPrice      float64 `json:"new_price"`
}

type PriceChange struct {
	ID                int64             `json:"id"`
	PriceChangeNumber string            `json:"price_change_number"`
	ChangeNumber      string            `json:"change_number"`
	ChangeDate        string            `json:"change_date"`
	EffectiveDate     string            `json:"effective_date"`
	Status            string            `json:"status"`
	Notes             string            `json:"notes"`
	CreatedBy         *int64            `json:"created_by"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	Items             []PriceChangeItem `json:"items"`
}

type PriceChangeRepository struct {
	db *sql.DB
}

func NewPriceChangeRepository(db *sql.DB) *PriceChangeRepository {
	return &PriceChangeRepository{db: db}
}

func (r *PriceChangeRepository) ListAll(ctx context.Context) ([]PriceChange, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PriceChange, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, change_number, effective_date::text, status, notes, created_by, created_at, updated_at
			FROM price_changes WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PriceChange, 0)
		for rows.Next() {
			var pc PriceChange
			if err := rows.Scan(&pc.ID, &pc.ChangeNumber, &pc.EffectiveDate, &pc.Status, &pc.Notes, &pc.CreatedBy, &pc.CreatedAt, &pc.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, pc.ID)
			pc.Items = items
			if pc.Items == nil {
				pc.Items = make([]PriceChangeItem, 0)
			}
			list = append(list, pc)
		}
		return list, nil
	})
}

func (r *PriceChangeRepository) GetByID(ctx context.Context, id int64) (*PriceChange, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*PriceChange, error) {
		pc := &PriceChange{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, change_number, effective_date::text, status, notes, created_by, created_at, updated_at
			FROM price_changes WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&pc.ID, &pc.ChangeNumber, &pc.EffectiveDate, &pc.Status, &pc.Notes, &pc.CreatedBy, &pc.CreatedAt, &pc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			pc.Items = items
		}
		return pc, nil
	})
}

func (r *PriceChangeRepository) GetItems(ctx context.Context, priceChangeID int64) ([]PriceChangeItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PriceChangeItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, price_change_id, product_id, product_name, old_price, new_price FROM price_change_items WHERE price_change_id=$1 AND tenant_id=$2 ORDER BY id", priceChangeID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]PriceChangeItem, 0)
		for rows.Next() {
			var i PriceChangeItem
			if err := rows.Scan(&i.ID, &i.PriceChangeID, &i.ProductID, &i.ProductName, &i.OldPrice, &i.NewPrice); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *PriceChangeRepository) Create(ctx context.Context, pc *PriceChange) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO price_changes (change_number, effective_date, notes, created_by, tenant_id) VALUES ($1,COALESCE(NULLIF($2,'')::date, CURRENT_DATE),$3,$4,$5) RETURNING id",
			pc.ChangeNumber, pc.EffectiveDate, pc.Notes, pc.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range pc.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO price_change_items (price_change_id, product_id, product_name, old_price, new_price, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.OldPrice, item.NewPrice, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *PriceChangeRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "UPDATE price_changes SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PriceChangeRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM price_changes WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *PriceChangeRepository) GeneratePriceChangeNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM price_changes WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("PC-%05d", count+1), nil
	})
}
