package usecase

import (
	"context"
	"errors"
	"math"
	"pos-system/internal/repository/postgres"
)

type JournalEntryUsecase struct {
	repo *postgres.JournalEntryRepository
}

func NewJournalEntryUsecase(repo *postgres.JournalEntryRepository) *JournalEntryUsecase {
	return &JournalEntryUsecase{repo: repo}
}

func (u *JournalEntryUsecase) ListAll(ctx context.Context) ([]postgres.JournalEntry, error) {
	return u.repo.ListAll(ctx)
}

func (u *JournalEntryUsecase) GetByID(ctx context.Context, id int64) (*postgres.JournalEntry, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *JournalEntryUsecase) Create(ctx context.Context, e *postgres.JournalEntry, items []postgres.JournalEntryItem) error {
	if len(items) < 2 {
		return errors.New("jurnal minimal memiliki 2 baris")
	}
	if e.Reference == "" {
		e.Reference = "Umum"
	}
	var debit, credit float64
	for _, it := range items {
		if it.AccountID <= 0 {
			return errors.New("setiap baris harus memilih akun")
		}
		if it.Debit < 0 || it.Credit < 0 {
			return errors.New("nilai debit/kredit tidak boleh negatif")
		}
		if it.Debit == 0 && it.Credit == 0 {
			return errors.New("setiap baris harus memiliki nilai debit atau kredit")
		}
		debit += it.Debit
		credit += it.Credit
	}
	if math.Abs(debit-credit) > 0.01 {
		return errors.New("total debit dan kredit harus seimbang")
	}
	number, err := u.repo.GenerateNumber(ctx)
	if err != nil {
		return err
	}
	e.EntryNumber = number
	e.Status = "posted"
	return u.repo.Create(ctx, e, items)
}

func (u *JournalEntryUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
