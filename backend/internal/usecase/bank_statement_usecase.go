package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
	"time"
)

type BankStatementUsecase struct {
	repo               *postgres.BankStatementRepository
	reconciliationRepo *postgres.BankReconciliationRepository
}

func NewBankStatementUsecase(repo *postgres.BankStatementRepository, reconciliationRepo *postgres.BankReconciliationRepository) *BankStatementUsecase {
	return &BankStatementUsecase{repo: repo, reconciliationRepo: reconciliationRepo}
}

func (u *BankStatementUsecase) ListImports(ctx context.Context) ([]postgres.BankStatementImport, error) {
	return u.repo.ListImports(ctx)
}

func (u *BankStatementUsecase) ImportLines(ctx context.Context, b *postgres.BankStatementImport) error {
	if b.AccountName == "" {
		return errors.New("nama rekening wajib diisi")
	}
	if len(b.Lines) == 0 {
		return errors.New("minimal harus ada 1 baris mutasi")
	}
	if b.Source == "" {
		b.Source = "manual"
	}
	if b.StatementDate == "" {
		b.StatementDate = time.Now().Format("2006-01-02")
	}
	b.TotalRows = len(b.Lines)
	for i := range b.Lines {
		if b.Lines[i].Type != "credit" && b.Lines[i].Type != "kredit" {
			b.Lines[i].Type = "debit"
		}
		if b.Lines[i].TransactionDate == "" {
			b.Lines[i].TransactionDate = b.StatementDate
		}
	}
	number, err := u.repo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	b.ImportNumber = number
	b.Status = "imported"
	id, err := u.repo.Create(ctx, b)
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (u *BankStatementUsecase) Reconcile(ctx context.Context, importID int64, bookBalance float64, createdBy *int64) error {
	imp, err := u.repo.GetByID(ctx, importID)
	if err != nil {
		return errors.New("data impor mutasi tidak ditemukan")
	}
	if len(imp.Lines) == 0 {
		return errors.New("tidak ada baris mutasi untuk direkonsiliasi")
	}

	statementBalance := bookBalance
	for _, l := range imp.Lines {
		if l.Type == "credit" || l.Type == "kredit" {
			statementBalance += l.Amount
		} else {
			statementBalance -= l.Amount
		}
	}

	number, err := u.reconciliationRepo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	br := &postgres.BankReconciliation{
		ReconciliationNumber: number,
		AccountName:          imp.AccountName,
		AccountType:          imp.AccountType,
		PeriodStart:          imp.StatementDate,
		PeriodEnd:            imp.StatementDate,
		BookBalance:          bookBalance,
		StatementBalance:     statementBalance,
		Difference:           statementBalance - bookBalance,
		Notes:                "Dibuat otomatis dari SmartLink e-Banking (" + imp.ImportNumber + ")",
		Status:               "completed",
		CreatedBy:            createdBy,
	}
	for _, l := range imp.Lines {
		br.Items = append(br.Items, postgres.BankReconciliationItem{
			Description:     l.Description,
			TransactionDate: l.TransactionDate,
			Amount:          l.Amount,
			Type:            l.Type,
			IsMatched:       true,
		})
	}

	reconID, err := u.reconciliationRepo.Create(ctx, br)
	if err != nil {
		return err
	}
	for _, l := range imp.Lines {
		_ = u.repo.MarkLineReconciled(ctx, l.ID, reconID)
	}
	return nil
}

func (u *BankStatementUsecase) DeleteImport(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
