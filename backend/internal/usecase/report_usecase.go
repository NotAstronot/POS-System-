package usecase

import (
	"context"
	"pos-system/internal/repository/postgres"
)

var reportPeriods = map[string]bool{
	"today":      true,
	"yesterday":  true,
	"this_week":  true,
	"last_week":  true,
	"this_month": true,
	"last_month": true,
	"this_year":  true,
	"last_year":  true,
}

type ReportUsecase struct {
	repo *postgres.ReportRepository
}

func NewReportUsecase(repo *postgres.ReportRepository) *ReportUsecase {
	return &ReportUsecase{repo: repo}
}

func (u *ReportUsecase) Outlets(ctx context.Context) ([]postgres.Outlet, error) {
	return u.repo.Outlets(ctx)
}

func (u *ReportUsecase) Summary(ctx context.Context, period string, outletID *int64) (*postgres.ReportSummary, error) {
	if !reportPeriods[period] {
		period = "this_month"
	}
	return u.repo.Summary(ctx, period, outletID)
}

func (u *ReportUsecase) PaymentMethods(ctx context.Context, period string, outletID *int64) ([]postgres.PaymentMethodStat, error) {
	if !reportPeriods[period] {
		period = "this_month"
	}
	return u.repo.PaymentMethods(ctx, period, outletID)
}
