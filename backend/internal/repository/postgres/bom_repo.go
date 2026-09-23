package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BomItem struct {
	ID               int64   `json:"id"`
	BomID            int64   `json:"bom_id"`
	ProductID        int64   `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         float64 `json:"quantity"`
	QuantityRequired float64 `json:"quantity_required"`
	Unit             string  `json:"unit"`
	CostPerUnit      float64 `json:"cost_per_unit"`
	Subtotal         float64 `json:"subtotal"`
	EstimatedCost    float64 `json:"estimated_cost"`
}

type BillOfMaterialsItem struct {
	ID               int64   `json:"id"`
	BomID            int64   `json:"bom_id"`
	ProductID        int64   `json:"product_id"`
	ProductName      string  `json:"product_name"`
	QuantityRequired float64 `json:"quantity_required"`
	Unit             string  `json:"unit"`
	CostPerUnit      float64 `json:"cost_per_unit"`
	EstimatedCost    float64 `json:"estimated_cost"`
}

type BillOfMaterials struct {
	ID             int64     `json:"id"`
	BomNumber      string    `json:"bom_number"`
	ProductID      int64     `json:"product_id"`
	ProductName    string    `json:"product_name"`
	Description    string    `json:"description"`
	QuantityOutput float64   `json:"quantity_output"`
	Unit           string    `json:"unit"`
	OverheadCost   float64   `json:"overhead_cost"`
	CostPerUnit    float64   `json:"cost_per_unit"`
	Status         string    `json:"status"`
	IsActive       bool      `json:"is_active"`
	Notes          string    `json:"notes"`
	CreatedBy      *int64    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Items          []BomItem `json:"items"`
}

type Bom struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   *int64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Items       []BomItem `json:"items"`
}

type BomRepository struct {
	db *sql.DB
}

func NewBomRepository(db *sql.DB) *BomRepository {
	return &BomRepository{db: db}
}

func (r *BomRepository) ListAll(ctx context.Context) ([]Bom, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Bom, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT b.id, b.product_id, COALESCE(p.name,''), b.description, b.is_active, b.created_by, b.created_at, b.updated_at
			FROM boms b
			LEFT JOIN products p ON p.id = b.product_id AND p.tenant_id = b.tenant_id
			WHERE b.tenant_id=$1
			ORDER BY b.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]Bom, 0)
		for rows.Next() {
			var b Bom
			if err := rows.Scan(&b.ID, &b.ProductID, &b.ProductName, &b.Description, &b.IsActive, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, b.ID)
			b.Items = items
			if b.Items == nil {
				b.Items = make([]BomItem, 0)
			}
			list = append(list, b)
		}
		return list, nil
	})
}

func (r *BomRepository) GetByID(ctx context.Context, id int64) (*BillOfMaterials, error) {
	return r.GetBillOfMaterialsByID(ctx, id)
}

func (r *BomRepository) GetItems(ctx context.Context, bomID int64) ([]BomItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BomItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, bom_id, product_id, product_name, quantity, unit FROM bom_items WHERE bom_id=$1 AND tenant_id=$2 ORDER BY id", bomID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BomItem, 0)
		for rows.Next() {
			var i BomItem
			if err := rows.Scan(&i.ID, &i.BomID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *BomRepository) Create(ctx context.Context, b *BillOfMaterials) (int64, error) {
	return r.CreateBillOfMaterials(ctx, b)
}

func (r *BomRepository) Update(ctx context.Context, b *BillOfMaterials) error {
	return r.UpdateBillOfMaterials(ctx, b)
}

func (r *BomRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM boms WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *BomRepository) List(ctx context.Context) ([]BillOfMaterials, error) {
	return r.ListBillOfMaterials(ctx)
}

func (r *BomRepository) ListBillOfMaterials(ctx context.Context) ([]BillOfMaterials, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BillOfMaterials, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT b.id, b.bom_number, b.product_id, COALESCE(p.name,''), COALESCE(b.description,''), b.quantity_output, COALESCE(b.unit,'pcs'), b.overhead_cost, b.cost_per_unit, b.status, b.created_by, b.created_at, b.updated_at
			FROM bill_of_materials b
			LEFT JOIN products p ON p.id = b.product_id AND p.tenant_id = b.tenant_id
			WHERE b.tenant_id=$1
			ORDER BY b.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BillOfMaterials, 0)
		for rows.Next() {
			var b BillOfMaterials
			if err := rows.Scan(&b.ID, &b.BomNumber, &b.ProductID, &b.ProductName, &b.Description, &b.QuantityOutput, &b.Unit, &b.OverheadCost, &b.CostPerUnit, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetBillOfMaterialsItems(ctx, b.ID)
			b.Items = items
			if b.Items == nil {
				b.Items = make([]BomItem, 0)
			}
			list = append(list, b)
		}
		return list, nil
	})
}

func (r *BomRepository) GetBillOfMaterialsByID(ctx context.Context, id int64) (*BillOfMaterials, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*BillOfMaterials, error) {
		b := &BillOfMaterials{}
		err := tx.QueryRowContext(ctx, `
			SELECT b.id, b.bom_number, b.product_id, COALESCE(p.name,''), COALESCE(b.description,''), b.quantity_output, COALESCE(b.unit,'pcs'), b.overhead_cost, b.cost_per_unit, b.status, b.created_by, b.created_at, b.updated_at
			FROM bill_of_materials b
			LEFT JOIN products p ON p.id = b.product_id AND p.tenant_id = b.tenant_id
			WHERE b.id=$1 AND b.tenant_id=$2`, id, tenantID).
			Scan(&b.ID, &b.BomNumber, &b.ProductID, &b.ProductName, &b.Description, &b.QuantityOutput, &b.Unit, &b.OverheadCost, &b.CostPerUnit, &b.Status, &b.CreatedBy, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetBillOfMaterialsItems(ctx, id)
		if err == nil {
			b.Items = items
		}
		return b, nil
	})
}

func (r *BomRepository) GetBillOfMaterialsItems(ctx context.Context, bomID int64) ([]BomItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]BomItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, bom_id, product_id, product_name, quantity_required, unit, cost_per_unit, estimated_cost FROM bill_of_material_items WHERE bom_id=$1 AND tenant_id=$2 ORDER BY id", bomID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]BomItem, 0)
		for rows.Next() {
			var i BomItem
			if err := rows.Scan(&i.ID, &i.BomID, &i.ProductID, &i.ProductName, &i.QuantityRequired, &i.Unit, &i.CostPerUnit, &i.EstimatedCost); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *BomRepository) CreateBillOfMaterials(ctx context.Context, b *BillOfMaterials) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO bill_of_materials (bom_number, product_id, description, quantity_output, unit, overhead_cost, cost_per_unit, status, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id",
			b.BomNumber, b.ProductID, b.Description, b.QuantityOutput, b.Unit, b.OverheadCost, b.CostPerUnit, b.Status, b.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range b.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO bill_of_material_items (bom_id, product_id, product_name, quantity_required, unit, cost_per_unit, estimated_cost, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
				id, item.ProductID, item.ProductName, item.QuantityRequired, item.Unit, item.CostPerUnit, item.EstimatedCost, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *BomRepository) UpdateBillOfMaterials(ctx context.Context, b *BillOfMaterials) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		if _, err := tx.ExecContext(ctx, "UPDATE bill_of_materials SET description=$1, quantity_output=$2, unit=$3, overhead_cost=$4, cost_per_unit=$5, status=$6, updated_at=NOW() WHERE id=$7 AND tenant_id=$8", b.Description, b.QuantityOutput, b.Unit, b.OverheadCost, b.CostPerUnit, b.Status, b.ID, tenantID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "DELETE FROM bill_of_material_items WHERE bom_id=$1 AND tenant_id=$2", b.ID, tenantID); err != nil {
			return err
		}
		for _, item := range b.Items {
			if _, err := tx.ExecContext(ctx,
				"INSERT INTO bill_of_material_items (bom_id, product_id, product_name, quantity_required, unit, cost_per_unit, estimated_cost, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
				b.ID, item.ProductID, item.ProductName, item.QuantityRequired, item.Unit, item.CostPerUnit, item.EstimatedCost, tenantID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *BomRepository) GenerateNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bill_of_materials WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("BOM-%05d", count+1), nil
	})
}
