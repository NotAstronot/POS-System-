package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type WorkOrderItem struct {
	ID               int64   `json:"id"`
	WorkOrderID      int64   `json:"work_order_id"`
	ProductID        int64   `json:"product_id"`
	ProductName      string  `json:"product_name"`
	Quantity         float64 `json:"quantity"`
	QuantityRequired float64 `json:"quantity_required"`
	Unit             string  `json:"unit"`
	CostPerUnit      float64 `json:"cost_per_unit"`
	Subtotal         float64 `json:"subtotal"`
}

type WorkOrder struct {
	ID              int64           `json:"id"`
	WorkOrderNumber string          `json:"work_order_number"`
	OrderNumber     string          `json:"order_number"`
	BomID           int64           `json:"bom_id"`
	BomNumber       string          `json:"bom_number"`
	Quantity        float64         `json:"quantity"`
	ProductID       int64           `json:"product_id"`
	ProductName     string          `json:"product_name"`
	TargetQty       float64         `json:"target_qty"`
	CompletedQty    float64         `json:"completed_qty"`
	ActualCost      float64         `json:"actual_cost"`
	WarehouseID     *int64          `json:"warehouse_id"`
	Status          string          `json:"status"`
	PlannedDate     string          `json:"planned_date"`
	ScheduledDate   string          `json:"scheduled_date"`
	CompletedDate   string          `json:"completed_date"`
	Notes           string          `json:"notes"`
	CreatedBy       *int64          `json:"created_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Items           []WorkOrderItem `json:"items"`
}

type WorkOrderRepository struct {
	db *sql.DB
}

func NewWorkOrderRepository(db *sql.DB) *WorkOrderRepository {
	return &WorkOrderRepository{db: db}
}

func (r *WorkOrderRepository) List(ctx context.Context) ([]WorkOrder, error) {
	return r.ListAll(ctx)
}

func (r *WorkOrderRepository) ListAll(ctx context.Context) ([]WorkOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]WorkOrder, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT wo.id, wo.order_number, wo.product_id, COALESCE(p.name,''), wo.target_qty, wo.completed_qty, wo.warehouse_id, wo.status, wo.planned_date::text, COALESCE(wo.completed_date::text,''), wo.notes, wo.created_by, wo.created_at, wo.updated_at
			FROM work_orders wo
			LEFT JOIN products p ON p.id = wo.product_id AND p.tenant_id = wo.tenant_id
			WHERE wo.tenant_id=$1
			ORDER BY wo.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]WorkOrder, 0)
		for rows.Next() {
			var wo WorkOrder
			if err := rows.Scan(&wo.ID, &wo.OrderNumber, &wo.ProductID, &wo.ProductName, &wo.TargetQty, &wo.CompletedQty, &wo.WarehouseID, &wo.Status, &wo.PlannedDate, &wo.CompletedDate, &wo.Notes, &wo.CreatedBy, &wo.CreatedAt, &wo.UpdatedAt); err != nil {
				return nil, err
			}
			items, _ := r.GetItems(ctx, wo.ID)
			wo.Items = items
			if wo.Items == nil {
				wo.Items = make([]WorkOrderItem, 0)
			}
			list = append(list, wo)
		}
		return list, nil
	})
}

func (r *WorkOrderRepository) GetByID(ctx context.Context, id int64) (*WorkOrder, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*WorkOrder, error) {
		wo := &WorkOrder{}
		err := tx.QueryRowContext(ctx, `
			SELECT wo.id, wo.order_number, wo.product_id, COALESCE(p.name,''), wo.target_qty, wo.completed_qty, wo.warehouse_id, wo.status, wo.planned_date::text, COALESCE(wo.completed_date::text,''), wo.notes, wo.created_by, wo.created_at, wo.updated_at
			FROM work_orders wo
			LEFT JOIN products p ON p.id = wo.product_id AND p.tenant_id = wo.tenant_id
			WHERE wo.id=$1 AND wo.tenant_id=$2`, id, tenantID).
			Scan(&wo.ID, &wo.OrderNumber, &wo.ProductID, &wo.ProductName, &wo.TargetQty, &wo.CompletedQty, &wo.WarehouseID, &wo.Status, &wo.PlannedDate, &wo.CompletedDate, &wo.Notes, &wo.CreatedBy, &wo.CreatedAt, &wo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items, err := r.GetItems(ctx, id)
		if err == nil {
			wo.Items = items
		}
		return wo, nil
	})
}

func (r *WorkOrderRepository) GetItems(ctx context.Context, workOrderID int64) ([]WorkOrderItem, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]WorkOrderItem, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, work_order_id, product_id, product_name, quantity, unit FROM work_order_items WHERE work_order_id=$1 AND tenant_id=$2 ORDER BY id", workOrderID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]WorkOrderItem, 0)
		for rows.Next() {
			var i WorkOrderItem
			if err := rows.Scan(&i.ID, &i.WorkOrderID, &i.ProductID, &i.ProductName, &i.Quantity, &i.Unit); err != nil {
				return nil, err
			}
			list = append(list, i)
		}
		return list, nil
	})
}

func (r *WorkOrderRepository) Create(ctx context.Context, wo *WorkOrder) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO work_orders (order_number, product_id, target_qty, warehouse_id, planned_date, notes, created_by, tenant_id) VALUES ($1,$2,$3,$4,COALESCE(NULLIF($5,'')::date, CURRENT_DATE),$6,$7,$8) RETURNING id",
			wo.OrderNumber, wo.ProductID, wo.TargetQty, wo.WarehouseID, wo.PlannedDate, wo.Notes, wo.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, item := range wo.Items {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO work_order_items (work_order_id, product_id, product_name, quantity, unit, tenant_id) VALUES ($1,$2,COALESCE(NULLIF($3,''), (SELECT name FROM products WHERE id=$2 AND tenant_id=$6)),$4,$5,$6)",
				id, item.ProductID, item.ProductName, item.Quantity, item.Unit, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *WorkOrderRepository) UpdateStatus(ctx context.Context, id int64, status string, completedQty float64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		var err error
		if status == "completed" {
			_, err = tx.ExecContext(ctx, "UPDATE work_orders SET status=$1, completed_qty=$2, completed_date=CURRENT_DATE, updated_at=NOW() WHERE id=$3 AND tenant_id=$4", status, completedQty, id, tenantID)
		} else {
			_, err = tx.ExecContext(ctx, "UPDATE work_orders SET status=$1, updated_at=NOW() WHERE id=$2 AND tenant_id=$3", status, id, tenantID)
		}
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WorkOrderRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM work_orders WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *WorkOrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM work_orders WHERE tenant_id=$1", tenantID).Scan(&count)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("WO-%05d", count+1), nil
	})
}

func (r *WorkOrderRepository) GenerateNumber(ctx context.Context) (string, error) {
	return r.GenerateOrderNumber(ctx)
}
