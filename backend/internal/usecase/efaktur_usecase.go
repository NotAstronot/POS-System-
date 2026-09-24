package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pos-system/internal/repository/postgres"
	"strings"
)

const (
	ExportTypePPNKeluaran = "ppn_keluaran"
	ExportTypePPNMasukan  = "ppn_masukan"
	ExportTypePPh23       = "pph23"
	ExportTypePPh21       = "pph21"
)

type EFakturUsecase struct {
	repo *postgres.EFakturRepository
}

func NewEFakturUsecase(repo *postgres.EFakturRepository) *EFakturUsecase {
	return &EFakturUsecase{repo: repo}
}

func (u *EFakturUsecase) ListExports(ctx context.Context) ([]postgres.EFakturExport, error) {
	return u.repo.ListExports(ctx)
}

func (u *EFakturUsecase) DeleteExport(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *EFakturUsecase) Generate(ctx context.Context, exportType, period string, createdBy *int64) (*postgres.EFakturExport, error) {
	if exportType == "" {
		return nil, errors.New("jenis ekspor wajib diisi")
	}
	if period == "" || len(period) != 7 {
		return nil, errors.New("periode wajib diisi dengan format YYYY-MM")
	}

	lines, err := u.buildLines(ctx, exportType, period)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("tidak ada data pajak pada periode tersebut")
	}

	number, err := u.repo.GenerateNumber(ctx, exportType)
	if err != nil {
		return nil, err
	}

	e := &postgres.EFakturExport{
		ExportNumber: number,
		ExportType:   exportType,
		Period:       period,
		Status:       "generated",
		TotalRows:    len(lines),
		CreatedBy:    createdBy,
		Lines:        lines,
	}
	for _, l := range lines {
		e.TotalDPP += l.TaxableBase
		e.TotalTax += l.TaxAmount
	}

	filePath, err := u.writeCSV(e)
	if err == nil {
		e.FilePath = filePath
	}

	id, err := u.repo.Create(ctx, e)
	if err != nil {
		return nil, err
	}
	e.ID = id
	return e, nil
}

func (u *EFakturUsecase) buildLines(ctx context.Context, exportType, period string) ([]postgres.EFakturLine, error) {
	src, err := u.repo.FetchExportSource(ctx, exportType, period)
	if err != nil {
		return nil, err
	}
	lines := make([]postgres.EFakturLine, 0, len(src))

	switch exportType {
	case ExportTypePPNKeluaran, ExportTypePPNMasukan:
		for _, row := range src {
			dpp := row.Amount / 1.11
			date := row.LineDate
			lines = append(lines, postgres.EFakturLine{
				ReferenceNumber: row.ReferenceNumber, PartyName: row.PartyName, NPWP: row.NPWP,
				LineDate: &date, TaxableBase: dpp, TaxAmount: row.Amount - dpp, Status: "valid",
			})
		}

	case ExportTypePPh23:
		for _, row := range src {
			date := row.LineDate
			lines = append(lines, postgres.EFakturLine{
				ReferenceNumber: row.ReferenceNumber, PartyName: row.PartyName, NPWP: row.NPWP,
				LineDate: &date, TaxableBase: row.Amount, TaxAmount: row.Amount * 0.02, Status: "valid",
			})
		}

	case ExportTypePPh21:
		for _, row := range src {
			lines = append(lines, postgres.EFakturLine{
				ReferenceNumber: "PPh21-" + row.PayPeriod,
				PartyName:       row.PartyName,
				LineDate:        ptrString(period + "-01"),
				TaxableBase:     row.Amount,
				TaxAmount:       row.Amount * 0.05,
				Status:          "valid",
			})
		}

	default:
		return nil, errors.New("jenis ekspor tidak dikenal")
	}
	return lines, nil
}

func (u *EFakturUsecase) writeCSV(e *postgres.EFakturExport) (string, error) {
	dir := "exports"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s_%s_%s.csv", e.ExportNumber, e.ExportType, strings.ReplaceAll(e.Period, "-", ""))
	fullPath := filepath.Join(dir, fileName)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	switch e.ExportType {
	case ExportTypePPh21:
		_ = w.Write([]string{"No", "Nama Karyawan", "Periode", "Penghasilan Bruto", "PPh21 (estimasi 5%)"})
	default:
		_ = w.Write([]string{"No", "Nomor Referensi", "Nama", "NPWP", "Tanggal", "Dasar Pengenaan Pajak (DPP)", "Pajak"})
	}

	for i, l := range e.Lines {
		date := ""
		if l.LineDate != nil {
			date = *l.LineDate
		}
		_ = w.Write([]string{
			fmt.Sprintf("%d", i+1),
			l.ReferenceNumber,
			l.PartyName,
			l.NPWP,
			date,
			fmt.Sprintf("%.2f", l.TaxableBase),
			fmt.Sprintf("%.2f", l.TaxAmount),
		})
	}
	return "/" + fullPath, nil
}

func ptrString(s string) *string { return &s }
