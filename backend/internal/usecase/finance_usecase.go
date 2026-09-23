package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type BankTransferUsecase struct {
	repo      *postgres.BankTransferRepository
	cashRepo  *postgres.CashTransactionRepository
	stockRepo *postgres.StockRepository
}

func NewBankTransferUsecase(repo *postgres.BankTransferRepository, cashRepo *postgres.CashTransactionRepository, stockRepo *postgres.StockRepository) *BankTransferUsecase {
	return &BankTransferUsecase{repo: repo, cashRepo: cashRepo, stockRepo: stockRepo}
}

func (u *BankTransferUsecase) ListAll(ctx context.Context) ([]postgres.BankTransfer, error) {
	return u.repo.ListAll(ctx)
}

func (u *BankTransferUsecase) GetByID(ctx context.Context, id int64) (*postgres.BankTransfer, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *BankTransferUsecase) Create(ctx context.Context, bt *postgres.BankTransfer) error {
	if bt.FromAccountName == "" {
		return errors.New("rekening asal wajib diisi")
	}
	if bt.ToAccountName == "" {
		return errors.New("rekening tujuan wajib diisi")
	}
	if bt.Amount <= 0 {
		return errors.New("jumlah transfer harus lebih dari 0")
	}
	if bt.FromAccountName == bt.ToAccountName && bt.FromAccountType == bt.ToAccountType {
		return errors.New("rekening asal dan tujuan tidak boleh sama")
	}
	number, err := u.repo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	bt.TransferNumber = number
	bt.Status = "completed"
	id, err := u.repo.Create(ctx, bt)
	if err != nil {
		return err
	}
	bt.ID = id

	if u.cashRepo != nil {
		if bt.FromAccountType == "kas" {
			_, _ = u.cashRepo.Create(ctx, &postgres.CashTransaction{
				Type: "out", Amount: bt.Amount,
				Description: "Transfer ke rekening: " + bt.ToAccountName,
				Category:    "transfer_bank", CreatedBy: bt.CreatedBy,
			})
		}
		if bt.ToAccountType == "kas" {
			_, _ = u.cashRepo.Create(ctx, &postgres.CashTransaction{
				Type: "in", Amount: bt.Amount,
				Description: "Transfer dari rekening: " + bt.FromAccountName,
				Category:    "transfer_bank", CreatedBy: bt.CreatedBy,
			})
		}
	}
	return nil
}

func (u *BankTransferUsecase) Update(ctx context.Context, bt *postgres.BankTransfer) error {
	return u.repo.Update(ctx, bt)
}

func (u *BankTransferUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

type BankReconciliationUsecase struct {
	repo *postgres.BankReconciliationRepository
}

func NewBankReconciliationUsecase(repo *postgres.BankReconciliationRepository) *BankReconciliationUsecase {
	return &BankReconciliationUsecase{repo: repo}
}

func (u *BankReconciliationUsecase) ListAll(ctx context.Context) ([]postgres.BankReconciliation, error) {
	return u.repo.ListAll(ctx)
}

func (u *BankReconciliationUsecase) GetByID(ctx context.Context, id int64) (*postgres.BankReconciliation, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *BankReconciliationUsecase) Create(ctx context.Context, br *postgres.BankReconciliation) error {
	if br.AccountName == "" {
		return errors.New("nama rekening wajib diisi")
	}
	if br.PeriodStart == "" || br.PeriodEnd == "" {
		return errors.New("periode awal dan akhir wajib diisi")
	}
	br.Difference = br.StatementBalance - br.BookBalance
	number, err := u.repo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	br.ReconciliationNumber = number
	br.Status = "completed"
	id, err := u.repo.Create(ctx, br)
	if err != nil {
		return err
	}
	br.ID = id
	return nil
}

func (u *BankReconciliationUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
