package postgres

import (
	"context"
	"database/sql"
	"time"
)

type FixedAsset struct {
	ID                      int64     `json:"id"`
	AssetCode               string    `json:"asset_code"`
	Name                    string    `json:"name"`
	AssetType               string    `json:"asset_type"`
	PurchaseDate            string    `json:"purchase_date"`
	Cost                    float64   `json:"cost"`
	SalvageValue            float64   `json:"salvage_value"`
	UsefulLifeMonths        int       `json:"useful_life_months"`
	DepreciationMethod      string    `json:"depreciation_method"`
	AccumulatedDepreciation float64   `json:"accumulated_depreciation"`
	DepreciationAccountID   int64     `json:"depreciation_account_id"`
	DepreciationAccountCode string    `json:"depreciation_account_code"`
	DepreciationAccountName string    `json:"depreciation_account_name"`
	AccumulatedAccountID    int64     `json:"accumulated_account_id"`
	AccumulatedAccountCode  string    `json:"accumulated_account_code"`
	AccumulatedAccountName  string    `json:"accumulated_account_name"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type FixedAssetRepository struct {
	db *sql.DB
}

func NewFixedAssetRepository(db *sql.DB) *FixedAssetRepository {
	return &FixedAssetRepository{db: db}
}

func (r *FixedAssetRepository) ListAll(ctx context.Context) ([]FixedAsset, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]FixedAsset, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT f.id, f.asset_code, f.name, f.asset_type, f.purchase_date::text, f.cost, f.salvage_value,
				f.useful_life_months, f.depreciation_method, f.accumulated_depreciation,
				f.depreciation_account_id, COALESCE(d.code,''), COALESCE(d.name,''),
				f.accumulated_account_id, COALESCE(a.code,''), COALESCE(a.name,''),
				f.is_active, f.created_at, f.updated_at
			FROM fixed_assets f
			LEFT JOIN chart_of_accounts d ON d.id = f.depreciation_account_id
			LEFT JOIN chart_of_accounts a ON a.id = f.accumulated_account_id
			WHERE f.tenant_id=$1
			ORDER BY f.id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]FixedAsset, 0)
		for rows.Next() {
			var f FixedAsset
			if err := rows.Scan(&f.ID, &f.AssetCode, &f.Name, &f.AssetType, &f.PurchaseDate, &f.Cost, &f.SalvageValue,
				&f.UsefulLifeMonths, &f.DepreciationMethod, &f.AccumulatedDepreciation,
				&f.DepreciationAccountID, &f.DepreciationAccountCode, &f.DepreciationAccountName,
				&f.AccumulatedAccountID, &f.AccumulatedAccountCode, &f.AccumulatedAccountName,
				&f.IsActive, &f.CreatedAt, &f.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, f)
		}
		return list, nil
	})
}

func (r *FixedAssetRepository) GetByID(ctx context.Context, id int64) (*FixedAsset, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*FixedAsset, error) {
		f := &FixedAsset{}
		err := tx.QueryRowContext(ctx, `
			SELECT f.id, f.asset_code, f.name, f.asset_type, f.purchase_date::text, f.cost, f.salvage_value,
				f.useful_life_months, f.depreciation_method, f.accumulated_depreciation,
				f.depreciation_account_id, COALESCE(d.code,''), COALESCE(d.name,''),
				f.accumulated_account_id, COALESCE(a.code,''), COALESCE(a.name,''),
				f.is_active, f.created_at, f.updated_at
			FROM fixed_assets f
			LEFT JOIN chart_of_accounts d ON d.id = f.depreciation_account_id
			LEFT JOIN chart_of_accounts a ON a.id = f.accumulated_account_id
			WHERE f.id=$1 AND f.tenant_id=$2`, id, tenantID).
			Scan(&f.ID, &f.AssetCode, &f.Name, &f.AssetType, &f.PurchaseDate, &f.Cost, &f.SalvageValue,
				&f.UsefulLifeMonths, &f.DepreciationMethod, &f.AccumulatedDepreciation,
				&f.DepreciationAccountID, &f.DepreciationAccountCode, &f.DepreciationAccountName,
				&f.AccumulatedAccountID, &f.AccumulatedAccountCode, &f.AccumulatedAccountName,
				&f.IsActive, &f.CreatedAt, &f.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return f, nil
	})
}

func (r *FixedAssetRepository) Create(ctx context.Context, f *FixedAsset) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx, `
			INSERT INTO fixed_assets (asset_code, name, asset_type, purchase_date, cost, salvage_value, useful_life_months,
				depreciation_method, accumulated_depreciation, depreciation_account_id, accumulated_account_id, is_active, tenant_id)
			VALUES ($1,$2,$3, COALESCE(NULLIF($4,'')::date, CURRENT_DATE), $5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
			f.AssetCode, f.Name, f.AssetType, f.PurchaseDate, f.Cost, f.SalvageValue, f.UsefulLifeMonths,
			f.DepreciationMethod, f.AccumulatedDepreciation, f.DepreciationAccountID, f.AccumulatedAccountID, f.IsActive, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	})
}

func (r *FixedAssetRepository) Update(ctx context.Context, f *FixedAsset) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, `
			UPDATE fixed_assets SET asset_code=$1, name=$2, asset_type=$3, purchase_date=COALESCE(NULLIF($4,'')::date, CURRENT_DATE),
				cost=$5, salvage_value=$6, useful_life_months=$7, depreciation_method=$8, accumulated_depreciation=$9,
				depreciation_account_id=$10, accumulated_account_id=$11, is_active=$12, updated_at=NOW()
			WHERE id=$13 AND tenant_id=$14`,
			f.AssetCode, f.Name, f.AssetType, f.PurchaseDate, f.Cost, f.SalvageValue, f.UsefulLifeMonths,
			f.DepreciationMethod, f.AccumulatedDepreciation, f.DepreciationAccountID, f.AccumulatedAccountID, f.IsActive, f.ID, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *FixedAssetRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM fixed_assets WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
