package usecase

import (
	"context"
	"errors"
	"pos-system/internal/repository/postgres"
)

type FixedAssetUsecase struct {
	repo *postgres.FixedAssetRepository
}

func NewFixedAssetUsecase(repo *postgres.FixedAssetRepository) *FixedAssetUsecase {
	return &FixedAssetUsecase{repo: repo}
}

func (u *FixedAssetUsecase) ListAll(ctx context.Context) ([]postgres.FixedAsset, error) {
	return u.repo.ListAll(ctx)
}

func (u *FixedAssetUsecase) GetByID(ctx context.Context, id int64) (*postgres.FixedAsset, error) {
	return u.repo.GetByID(ctx, id)
}

func validateAsset(f *postgres.FixedAsset) error {
	switch f.AssetType {
	case "peralatan", "komputer", "kendaraan", "bangunan", "lainnya":
	default:
		return errors.New("tipe aset tidak valid")
	}
	switch f.DepreciationMethod {
	case "garis_lurus", "saldo_menurun":
	default:
		return errors.New("metode penyusutan tidak valid")
	}
	if f.Cost <= 0 {
		return errors.New("harga perolehan harus lebih dari 0")
	}
	if f.SalvageValue > f.Cost {
		return errors.New("nilai sisa tidak boleh melebihi harga perolehan")
	}
	if f.UsefulLifeMonths <= 0 {
		return errors.New("umur manfaat harus lebih dari 0 bulan")
	}
	if f.DepreciationAccountID <= 0 || f.AccumulatedAccountID <= 0 {
		return errors.New("pilih akun beban depresiasi dan akun akumulasi depresiasi")
	}
	return nil
}

func (u *FixedAssetUsecase) Create(ctx context.Context, f *postgres.FixedAsset) error {
	if err := validateAsset(f); err != nil {
		return err
	}
	id, err := u.repo.Create(ctx, f)
	if err != nil {
		return err
	}
	f.ID = id
	return nil
}

func (u *FixedAssetUsecase) Update(ctx context.Context, f *postgres.FixedAsset) error {
	if err := validateAsset(f); err != nil {
		return err
	}
	return u.repo.Update(ctx, f)
}

func (u *FixedAssetUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *FixedAssetUsecase) Schedule(ctx context.Context) ([]AssetDepreciationSchedule, error) {
	assets, err := u.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	schedule := make([]AssetDepreciationSchedule, 0, len(assets))
	for _, a := range assets {
		remaining := a.Cost - a.SalvageValue - a.AccumulatedDepreciation
		var monthly float64
		if a.UsefulLifeMonths > 0 {
			if a.DepreciationMethod == "saldo_menurun" {
				base := a.Cost - a.AccumulatedDepreciation
				if base > 0 {
					monthly = 2 * base / float64(a.UsefulLifeMonths)
				}
			} else {
				monthly = (a.Cost - a.SalvageValue) / float64(a.UsefulLifeMonths)
			}
		}
		if monthly < 0 {
			monthly = 0
		}
		first := monthly
		if first > remaining {
			first = remaining
		}
		if remaining < 0 {
			remaining = 0
		}
		monthsRemaining := 0.0
		if monthly > 0 {
			monthsRemaining = remaining / monthly
		}
		schedule = append(schedule, AssetDepreciationSchedule{
			FixedAsset:          a,
			MonthlyDepreciation: monthly,
			FirstDepreciation:   first,
			RemainingValue:      remaining,
			MonthsRemaining:     monthsRemaining,
		})
	}
	return schedule, nil
}

type AssetDepreciationSchedule struct {
	postgres.FixedAsset
	MonthlyDepreciation float64 `json:"monthly_depreciation"`
	FirstDepreciation   float64 `json:"first_depreciation"`
	RemainingValue      float64 `json:"remaining_value"`
	MonthsRemaining     float64 `json:"months_remaining"`
}
