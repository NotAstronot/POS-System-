package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"pos-system/internal/repository/postgres"
)

var validDiscounts = map[string]float64{
	"HEMAT10":  10,
	"DISKON15": 15,
	"PROMO20":  20,
	"MEMBER":   5,
}

type OrderUsecase struct {
	orderRepo   *postgres.OrderRepository
	revenueRepo *postgres.RevenueRepository
	shiftRepo   *postgres.ShiftRepository
}

func NewOrderUsecase(orderRepo *postgres.OrderRepository, revenueRepo *postgres.RevenueRepository, shiftRepo *postgres.ShiftRepository) *OrderUsecase {
	return &OrderUsecase{orderRepo: orderRepo, revenueRepo: revenueRepo, shiftRepo: shiftRepo}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, in postgres.CreateOrderInput) (*postgres.OrderResult, error) {
	if len(in.Items) == 0 {
		return nil, errors.New("keranjang masih kosong")
	}
	for _, p := range in.Payments {
		if p.Amount <= 0 {
			return nil, errors.New("jumlah pembayaran harus lebih dari 0")
		}
	}
	if in.OrderType == "" {
		in.OrderType = "dine_in"
	}
	if in.DiscountAmount <= 0 && in.DiscountCode != "" {
		if pct, ok := validDiscounts[strings.ToUpper(strings.TrimSpace(in.DiscountCode))]; ok {
			in.DiscountPercent = pct
		}
	}
	return u.orderRepo.CreateOrder(ctx, in)
}

type SyncOrderResult struct {
	ClientOrderID string `json:"client_order_id"`
	Success       bool   `json:"success"`
	AlreadySynced bool   `json:"already_synced"`
	OrderID       int64  `json:"order_id,omitempty"`
	OrderNumber   string `json:"order_number,omitempty"`
	Error         string `json:"error,omitempty"`
}

type SyncBatchResult struct {
	Total    int               `json:"total"`
	Synced   int               `json:"synced"`
	Failed   int               `json:"failed"`
	Results  []SyncOrderResult `json:"results"`
	SyncedAt string            `json:"synced_at"`
}

func (u *OrderUsecase) SyncOrders(ctx context.Context, batch []postgres.CreateOrderInput) *SyncBatchResult {
	res := &SyncBatchResult{Total: len(batch)}
	for _, in := range batch {
		if in.ClientOrderID == "" {
			in.ClientOrderID = fmt.Sprintf("LCL-%d", time.Now().UnixNano())
		}
		in.SkipShiftValidation = true
		if len(in.Items) == 0 {
			res.Failed++
			res.Results = append(res.Results, SyncOrderResult{ClientOrderID: in.ClientOrderID, Error: "order tidak memiliki item"})
			continue
		}
		result, err := u.CreateOrder(ctx, in)
		if err != nil {
			res.Failed++
			res.Results = append(res.Results, SyncOrderResult{ClientOrderID: in.ClientOrderID, Error: err.Error()})
			continue
		}
		res.Synced++
		res.Results = append(res.Results, SyncOrderResult{
			ClientOrderID: in.ClientOrderID,
			Success:       true,
			AlreadySynced: result.AlreadySynced,
			OrderID:       result.Order.ID,
			OrderNumber:   result.Order.OrderNumber,
		})
	}
	res.SyncedAt = time.Now().Format(time.RFC3339)
	return res
}

func (u *OrderUsecase) RecentOrders(ctx context.Context, limit int) ([]postgres.RecentOrder, error) {
	return u.orderRepo.RecentOrders(ctx, limit)
}

func (u *OrderUsecase) RevenueSummary(ctx context.Context, period string) (*postgres.RevenueSummary, error) {
	if period != "daily" && period != "weekly" && period != "monthly" {
		period = "daily"
	}
	return u.revenueRepo.Summary(ctx, period)
}

func (u *OrderUsecase) OpenShift(ctx context.Context, userID int64, openingBalance float64) (*postgres.Shift, error) {
	return u.shiftRepo.Create(ctx, userID, openingBalance)
}

func (u *OrderUsecase) CloseShift(ctx context.Context, id int64, closingBalance float64) error {
	return u.shiftRepo.Close(ctx, id, closingBalance)
}
