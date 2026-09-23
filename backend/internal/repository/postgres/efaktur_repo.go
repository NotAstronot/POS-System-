package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Efaktur struct {
	ID           int64     `json:"id"`
	FakturNo     string    `json:"faktur_no"`
	CustomerName string    `json:"customer_name"`
	Address      string    `json:"address"`
	TaxID        string    `json:"tax_id"`
	Date         string    `json:"date"`
	Status       string    `json:"status"`
	CreatedBy    *int64    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EFakturExport struct {
	ID           int64         `json:"id"`
	ExportNumber string        `json:"export_number"`
	ExportType   string        `json:"export_type"`
	Period       string        `json:"period"`
	Status       string        `json:"status"`
	TotalRows    int           `json:"total_rows"`
	TotalDPP     float64       `json:"total_dpp"`
	TotalTax     float64       `json:"total_tax"`
	FilePath     string        `json:"file_path"`
	CreatedBy    *int64        `json:"created_by"`
	CreatedAt    time.Time     `json:"created_at"`
	Lines        []EFakturLine `json:"lines"`
}

type EFakturLine struct {
	ID              int64   `json:"id"`
	ExportID        int64   `json:"export_id"`
	ReferenceNumber string  `json:"reference_number"`
	PartyName       string  `json:"party_name"`
	NPWP            string  `json:"npwp"`
	LineDate        *string `json:"line_date"`
	TaxableBase     float64 `json:"taxable_base"`
	TaxAmount       float64 `json:"tax_amount"`
	Status          string  `json:"status"`
}

type EFakturRepository struct {
	db *sql.DB
}

func NewEfakturRepository(db *sql.DB) *EFakturRepository {
	return &EFakturRepository{db: db}
}

func (r *EFakturRepository) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

func (r *EFakturRepository) ListAll(ctx context.Context) ([]Efaktur, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Efaktur, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, faktur_no, customer_name, address, tax_id, date::text, status, created_by, created_at, updated_at
			FROM efakturs WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]Efaktur, 0)
		for rows.Next() {
			var e Efaktur
			if err := rows.Scan(&e.ID, &e.FakturNo, &e.CustomerName, &e.Address, &e.TaxID, &e.Date, &e.Status, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
				return nil, err
			}
			list = append(list, e)
		}
		return list, nil
	})
}

func (r *EFakturRepository) ListExports(ctx context.Context) ([]EFakturExport, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]EFakturExport, error) {
		rows, err := tx.QueryContext(ctx, `
			SELECT id, export_number, export_type, period, status, total_rows, total_dpp, total_tax, COALESCE(file_path,''), created_by, created_at
			FROM efaktur_exports WHERE tenant_id=$1 ORDER BY id DESC`, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]EFakturExport, 0)
		for rows.Next() {
			var e EFakturExport
			if err := rows.Scan(&e.ID, &e.ExportNumber, &e.ExportType, &e.Period, &e.Status, &e.TotalRows, &e.TotalDPP, &e.TotalTax, &e.FilePath, &e.CreatedBy, &e.CreatedAt); err != nil {
				return nil, err
			}
			lines, _ := r.GetExportLines(ctx, e.ID)
			e.Lines = lines
			if e.Lines == nil {
				e.Lines = make([]EFakturLine, 0)
			}
			list = append(list, e)
		}
		return list, nil
	})
}

func (r *EFakturRepository) GetExportLines(ctx context.Context, exportID int64) ([]EFakturLine, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]EFakturLine, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, export_id, reference_number, party_name, npwp, line_date, taxable_base, tax_amount, status FROM efaktur_lines WHERE export_id=$1 AND tenant_id=$2 ORDER BY id", exportID, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		list := make([]EFakturLine, 0)
		for rows.Next() {
			var l EFakturLine
			if err := rows.Scan(&l.ID, &l.ExportID, &l.ReferenceNumber, &l.PartyName, &l.NPWP, &l.LineDate, &l.TaxableBase, &l.TaxAmount, &l.Status); err != nil {
				return nil, err
			}
			list = append(list, l)
		}
		return list, nil
	})
}

func (r *EFakturRepository) GetByID(ctx context.Context, id int64) (*Efaktur, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (*Efaktur, error) {
		e := &Efaktur{}
		err := tx.QueryRowContext(ctx, `
			SELECT id, faktur_no, customer_name, address, tax_id, date::text, status, created_by, created_at, updated_at
			FROM efakturs WHERE id=$1 AND tenant_id=$2`, id, tenantID).
			Scan(&e.ID, &e.FakturNo, &e.CustomerName, &e.Address, &e.TaxID, &e.Date, &e.Status, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return e, nil
	})
}

func (r *EFakturRepository) Create(ctx context.Context, e *EFakturExport) (int64, error) {
	return r.CreateExport(ctx, e)
}

func (r *EFakturRepository) CreateExport(ctx context.Context, e *EFakturExport) (int64, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (int64, error) {
		var id int64
		err := tx.QueryRowContext(ctx,
			"INSERT INTO efaktur_exports (export_number, export_type, period, status, total_rows, total_dpp, total_tax, file_path, created_by, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id",
			e.ExportNumber, e.ExportType, e.Period, e.Status, e.TotalRows, e.TotalDPP, e.TotalTax, e.FilePath, e.CreatedBy, tenantID).Scan(&id)
		if err != nil {
			return 0, err
		}
		for _, l := range e.Lines {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO efaktur_lines (export_id, reference_number, party_name, npwp, line_date, taxable_base, tax_amount, status, tenant_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
				id, l.ReferenceNumber, l.PartyName, l.NPWP, l.LineDate, l.TaxableBase, l.TaxAmount, l.Status, tenantID)
			if err != nil {
				return 0, err
			}
		}
		return id, nil
	})
}

func (r *EFakturRepository) GenerateNumber(ctx context.Context, exportType string) (string, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) (string, error) {
		var count int64
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM efaktur_exports WHERE tenant_id=$1 AND export_type=$2", tenantID, exportType).Scan(&count)
		if err != nil {
			return "", err
		}
		prefix := "EF"
		switch exportType {
		case "ppn_keluaran":
			prefix = "PPNK"
		case "ppn_masukan":
			prefix = "PPNM"
		case "pph23":
			prefix = "PPh23"
		case "pph21":
			prefix = "PPh21"
		}
		return fmt.Sprintf("%s-%05d", prefix, count+1), nil
	})
}

func (r *EFakturRepository) Delete(ctx context.Context, id int64) error {
	return withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		_, err := tx.ExecContext(ctx, "DELETE FROM efakturs WHERE id=$1 AND tenant_id=$2", id, tenantID)
		if err != nil {
			return err
		}
		return nil
	})
}
