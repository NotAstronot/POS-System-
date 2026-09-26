package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Outlet struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ReportKPIs struct {
	GrossSales     float64 `json:"gross_sales"`
	NetSales       float64 `json:"net_sales"`
	GrossProfit    float64 `json:"gross_profit"`
	Transactions   int     `json:"transactions"`
	AvgSalePerTx   float64 `json:"avg_sale_per_tx"`
	GrossMarginPct float64 `json:"gross_margin_pct"`
	ItemsSold      int     `json:"items_sold"`
	PrevGrossSales float64 `json:"prev_gross_sales"`
	ChangePct      float64 `json:"change_pct"`
}

type SalesPoint struct {
	Label        string  `json:"label"`
	GrossSales   float64 `json:"gross_sales"`
	NetSales     float64 `json:"net_sales"`
	Transactions int     `json:"transactions"`
}

type DayOfWeekSales struct {
	Day          string  `json:"day"`
	GrossSales   float64 `json:"gross_sales"`
	Transactions int     `json:"transactions"`
}

type HourlySales struct {
	Hour         string  `json:"hour"`
	GrossSales   float64 `json:"gross_sales"`
	Transactions int     `json:"transactions"`
}

type ReportItem struct {
	Name        string  `json:"name"`
	ItemsSold   int     `json:"items_sold"`
	GrossSales  float64 `json:"gross_sales"`
	NetSales    float64 `json:"net_sales"`
	GrossProfit float64 `json:"gross_profit"`
}

type CategoryStat struct {
	Category string  `json:"category"`
	Value    float64 `json:"value"`
	Pct      float64 `json:"pct"`
}

type TopItem struct {
	Name       string  `json:"name"`
	Qty        int     `json:"qty"`
	GrossSales float64 `json:"gross_sales"`
}

type CategoryTopItems struct {
	Category string    `json:"category"`
	Items    []TopItem `json:"items"`
}

type ReportSummary struct {
	Period             string             `json:"period"`
	RangeStart         string             `json:"range_start"`
	RangeEnd           string             `json:"range_end"`
	BucketMode         string             `json:"bucket_mode"`
	KPIs               ReportKPIs         `json:"kpis"`
	DailySales         []SalesPoint       `json:"daily_sales"`
	DayOfWeekSales     []DayOfWeekSales   `json:"day_of_week_sales"`
	HourlySales        []HourlySales      `json:"hourly_sales"`
	Items              []ReportItem       `json:"items"`
	CategoryByVolume   []CategoryStat     `json:"category_by_volume"`
	CategoryBySales    []CategoryStat     `json:"category_by_sales"`
	TopItemsByCategory []CategoryTopItems `json:"top_items_by_category"`
}

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func outletUUID(id int64) string {
	return fmt.Sprintf("b%035d", id)
}

func (r *ReportRepository) Outlets(ctx context.Context) ([]Outlet, error) {
	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]Outlet, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id, name FROM branches WHERE is_active = true AND tenant_id=$1 ORDER BY name", tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var list []Outlet
		for rows.Next() {
			var o Outlet
			if err := rows.Scan(&o.ID, &o.Name); err != nil {
				return nil, err
			}
			list = append(list, o)
		}
		return list, rows.Err()
	})
}

func (r *ReportRepository) Summary(ctx context.Context, period string, outletID *int64) (*ReportSummary, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)

	since, until, prevSince, prevUntil, bucketMode, ok := periodRange(period, now)
	if !ok {
		return nil, fmt.Errorf("periode tidak valid: %s", period)
	}

	res := &ReportSummary{
		Period:     period,
		RangeStart: since.Format("2006-01-02"),
		RangeEnd:   until.AddDate(0, 0, -1).Format("2006-01-02"),
		BucketMode: bucketMode,
	}

	var gross, net, tax, itemsSold float64
	var transactions int
	var cogs float64
	var prevGross float64

	err = withTenantTx(ctx, r.db, func(tx *sql.Tx, tenantID int64) error {
		where := "o.status='completed' AND o.created_at >= $1 AND o.created_at < $2 AND o.tenant_id=$3"
		args := []any{since, until, tenantID}
		argN := 3
		if outletID != nil {
			argN++
			where += fmt.Sprintf(" AND o.outlet_id = $%d", argN)
			args = append(args, outletUUID(*outletID))
		}

		err := tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(gross),0), COALESCE(SUM(gross-discount),0),
				COALESCE(SUM(total-gross+discount),0), COUNT(*), COALESCE(SUM(items),0)
			FROM (
				SELECT o.id, SUM(oi.subtotal) gross, o.discount_amount discount, o.total total, SUM(oi.quantity) items
				FROM orders o JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
				WHERE `+where+`
				GROUP BY o.id, o.discount_amount, o.total
			) t`, args...).Scan(&gross, &net, &tax, &transactions, &itemsSold)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		err = tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(oi.quantity * COALESCE(p.purchase_price,0)),0)
			FROM order_items oi
			JOIN orders o ON o.id = oi.order_id AND o.tenant_id = oi.tenant_id
			LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id = oi.tenant_id
			WHERE `+strings.Replace(where, "o.status", "o.status", 1), args...).Scan(&cogs)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		prevArgs := []any{prevSince, prevUntil, tenantID}
		if outletID != nil {
			prevArgs = append(prevArgs, outletUUID(*outletID))
		}
		prevWhere := "o.status='completed' AND o.created_at >= $1 AND o.created_at < $2 AND o.tenant_id=$3"
		if outletID != nil {
			prevWhere += fmt.Sprintf(" AND o.outlet_id = '%s'", outletUUID(*outletID))
		}
		err = tx.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(gross),0)
			FROM (
				SELECT o.id, SUM(oi.subtotal) gross
				FROM orders o JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
				WHERE `+prevWhere+`
				GROUP BY o.id
			) t`, prevArgs...).Scan(&prevGross)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		res.KPIs = ReportKPIs{
			GrossSales:     gross,
			NetSales:       net,
			GrossProfit:    net - cogs,
			Transactions:   transactions,
			ItemsSold:      int(itemsSold),
			PrevGrossSales: prevGross,
		}
		if transactions > 0 {
			res.KPIs.AvgSalePerTx = net / float64(transactions)
		}
		if net > 0 {
			res.KPIs.GrossMarginPct = (net - cogs) / net * 100
		}
		if prevGross > 0 {
			res.KPIs.ChangePct = (gross - prevGross) / prevGross * 100
		}

		if err := r.loadSalesBucketsTx(ctx, tx, res, bucketMode, since, until, where, args); err != nil {
			return err
		}
		if err := r.loadDayOfWeekTx(ctx, tx, res, since, until, where, args); err != nil {
			return err
		}
		if err := r.loadHourlyTx(ctx, tx, res, since, until, where, args); err != nil {
			return err
		}
		if err := r.loadItemsTx(ctx, tx, res, since, until, where, args); err != nil {
			return err
		}
		if err := r.loadCategoriesTx(ctx, tx, res, since, until, where, args); err != nil {
			return err
		}
		if err := r.loadTopItemsByCategoryTx(ctx, tx, res, since, until, where, args); err != nil {
			return err
		}
		return nil
	})

	return res, err
}

func outletWhereSQL(outletID *int64) string {
	if outletID == nil {
		return ""
	}
	return fmt.Sprintf(" AND o.outlet_id = '%s'", outletUUID(*outletID))
}

func periodRange(period string, now time.Time) (since, until, prevSince, prevUntil time.Time, bucketMode string, ok bool) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch period {
	case "today":
		since, until = today, today.AddDate(0, 0, 1)
		prevSince, prevUntil = today.AddDate(0, 0, -1), today
		bucketMode = "hourly"
	case "yesterday":
		since, until = today.AddDate(0, 0, -1), today
		prevSince, prevUntil = today.AddDate(0, 0, -2), today.AddDate(0, 0, -1)
		bucketMode = "hourly"
	case "this_week":
		offset := (int(now.Weekday()) + 6) % 7
		start := today.AddDate(0, 0, -offset)
		since, until = start, start.AddDate(0, 0, 7)
		prevSince, prevUntil = start.AddDate(0, 0, -7), start
		bucketMode = "daily"
	case "last_week":
		offset := (int(now.Weekday()) + 6) % 7
		start := today.AddDate(0, 0, -offset-7)
		since, until = start, start.AddDate(0, 0, 7)
		prevSince, prevUntil = start.AddDate(0, 0, -7), start
		bucketMode = "daily"
	case "this_month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		since, until = start, start.AddDate(0, 1, 0)
		prevSince, prevUntil = start.AddDate(0, -1, 0), start
		bucketMode = "daily"
	case "last_month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
		since, until = start, start.AddDate(0, 1, 0)
		prevSince, prevUntil = start.AddDate(0, -1, 0), start
		bucketMode = "daily"
	case "this_year":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		since, until = start, start.AddDate(1, 0, 0)
		prevSince, prevUntil = start.AddDate(-1, 0, 0), start
		bucketMode = "monthly"
	case "last_year":
		start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		since, until = start, start.AddDate(1, 0, 0)
		prevSince, prevUntil = start.AddDate(-1, 0, 0), start
		bucketMode = "monthly"
	default:
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, "", false
	}
	return since, until, prevSince, prevUntil, bucketMode, true
}

func (r *ReportRepository) loadSalesBucketsTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, bucketMode string, since, until time.Time, where string, args []any) error {
	var trunc string
	switch bucketMode {
	case "hourly":
		trunc = "EXTRACT(HOUR FROM o.created_at AT TIME ZONE 'Asia/Jakarta')"
	case "daily":
		trunc = "TO_CHAR(o.created_at AT TIME ZONE 'Asia/Jakarta','YYYY-MM-DD')"
	default:
		trunc = "TO_CHAR(o.created_at AT TIME ZONE 'Asia/Jakarta','YYYY-MM')"
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT `+trunc+` AS label, COUNT(*), SUM(gross), SUM(gross - discount)
		FROM (
			SELECT o.id, o.discount_amount,
				SUM(oi.subtotal) gross,
				o.created_at
			FROM orders o JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
			WHERE `+where+`
			GROUP BY o.id, o.discount_amount, o.created_at
		) t GROUP BY label`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	data := map[string]*SalesPoint{}
	var order []string
	for rows.Next() {
		var label string
		var sp SalesPoint
		if err := rows.Scan(&label, &sp.Transactions, &sp.GrossSales, &sp.NetSales); err != nil {
			return err
		}
		sp.Label = label
		if _, exists := data[label]; !exists {
			order = append(order, label)
		}
		data[label] = &sp
	}
	if err := rows.Err(); err != nil {
		return err
	}

	res.DailySales = []SalesPoint{}
	switch bucketMode {
	case "hourly":
		for h := 0; h < 24; h++ {
			label := fmt.Sprintf("%02d:00", h)
			if sp, ok := data[label]; ok {
				res.DailySales = append(res.DailySales, *sp)
			} else {
				res.DailySales = append(res.DailySales, SalesPoint{Label: label})
			}
		}
	case "daily":
		for d := since; d.Before(until); d = d.AddDate(0, 0, 1) {
			label := d.Format("2006-01-02")
			if sp, ok := data[label]; ok {
				res.DailySales = append(res.DailySales, *sp)
			} else {
				res.DailySales = append(res.DailySales, SalesPoint{Label: label})
			}
		}
	default:
		cur := since
		for cur.Before(until) {
			label := cur.Format("2006-01")
			if sp, ok := data[label]; ok {
				res.DailySales = append(res.DailySales, *sp)
			} else {
				res.DailySales = append(res.DailySales, SalesPoint{Label: label})
			}
			cur = cur.AddDate(0, 1, 0)
		}
	}
	return nil
}

func (r *ReportRepository) loadDayOfWeekTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, since, until time.Time, where string, args []any) error {
	rows, err := tx.QueryContext(ctx, `
		SELECT EXTRACT(ISODOW FROM o.created_at AT TIME ZONE 'Asia/Jakarta')::int, COUNT(*), SUM(gross)
		FROM (
			SELECT o.id, SUM(oi.subtotal) gross, o.created_at
			FROM orders o JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
			WHERE `+where+`
			GROUP BY o.id, o.created_at
		) t GROUP BY 1`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	byDow := map[int]*DayOfWeekSales{}
	for rows.Next() {
		var dow int
		var d DayOfWeekSales
		if err := rows.Scan(&dow, &d.Transactions, &d.GrossSales); err != nil {
			return err
		}
		byDow[dow] = &d
	}
	if err := rows.Err(); err != nil {
		return err
	}

	names := []string{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"}
	res.DayOfWeekSales = []DayOfWeekSales{}
	for i, name := range names {
		dow := i + 1
		if d, ok := byDow[dow]; ok {
			d.Day = name
			res.DayOfWeekSales = append(res.DayOfWeekSales, *d)
		} else {
			res.DayOfWeekSales = append(res.DayOfWeekSales, DayOfWeekSales{Day: name})
		}
	}
	return nil
}

func (r *ReportRepository) loadHourlyTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, since, until time.Time, where string, args []any) error {
	rows, err := tx.QueryContext(ctx, `
		SELECT EXTRACT(HOUR FROM o.created_at AT TIME ZONE 'Asia/Jakarta')::int, COUNT(*), SUM(gross)
		FROM (
			SELECT o.id, SUM(oi.subtotal) gross, o.created_at
			FROM orders o JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
			WHERE `+where+`
			GROUP BY o.id, o.created_at
		) t GROUP BY 1`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	byHour := map[int]*HourlySales{}
	for rows.Next() {
		var h int
		var s HourlySales
		if err := rows.Scan(&h, &s.Transactions, &s.GrossSales); err != nil {
			return err
		}
		byHour[h] = &s
	}
	if err := rows.Err(); err != nil {
		return err
	}

	res.HourlySales = []HourlySales{}
	for h := 0; h < 24; h++ {
		if s, ok := byHour[h]; ok {
			s.Hour = fmt.Sprintf("%02d:00", h)
			res.HourlySales = append(res.HourlySales, *s)
		} else {
			res.HourlySales = append(res.HourlySales, HourlySales{Hour: fmt.Sprintf("%02d:00", h)})
		}
	}
	return nil
}

func (r *ReportRepository) loadItemsTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, since, until time.Time, where string, args []any) error {
	rows, err := tx.QueryContext(ctx, `
		SELECT COALESCE(p.name,'Produk'), SUM(oi.quantity),
			SUM(oi.subtotal),
			SUM(oi.subtotal - o.discount_amount * (oi.subtotal / NULLIF(og.gross,0))),
			SUM(oi.quantity * COALESCE(p.purchase_price,0))
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
		LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id=o.tenant_id
		LEFT JOIN (
			SELECT o2.id, SUM(oi2.subtotal) gross
			FROM orders o2 JOIN order_items oi2 ON oi2.order_id = o2.id AND oi2.tenant_id=o2.tenant_id
			GROUP BY o2.id
		) og ON og.id = o.id
		WHERE `+where+`
		GROUP BY p.id, p.name
		ORDER BY SUM(oi.subtotal) DESC LIMIT 30`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	res.Items = []ReportItem{}
	for rows.Next() {
		var it ReportItem
		var cogs float64
		if err := rows.Scan(&it.Name, &it.ItemsSold, &it.GrossSales, &it.NetSales, &cogs); err != nil {
			return err
		}
		it.GrossProfit = it.NetSales - cogs
		res.Items = append(res.Items, it)
	}
	return rows.Err()
}

func (r *ReportRepository) loadCategoriesTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, since, until time.Time, where string, args []any) error {
	rows, err := tx.QueryContext(ctx, `
		SELECT COALESCE(c.name,'Tanpa Kategori'), SUM(oi.quantity), SUM(oi.subtotal)
		FROM orders o
		JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
		LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id=o.tenant_id
		LEFT JOIN categories c ON c.id = p.category_id AND c.tenant_id=o.tenant_id
		WHERE `+where+`
		GROUP BY 1`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var volumeTotal, salesTotal float64
	volume := map[string]float64{}
	sales := map[string]float64{}
	var order []string
	for rows.Next() {
		var cat string
		var v, s float64
		if err := rows.Scan(&cat, &v, &s); err != nil {
			return err
		}
		if _, exists := volume[cat]; !exists {
			order = append(order, cat)
		}
		volume[cat] += v
		sales[cat] += s
		volumeTotal += v
		salesTotal += s
	}
	if err := rows.Err(); err != nil {
		return err
	}

	res.CategoryByVolume = []CategoryStat{}
	res.CategoryBySales = []CategoryStat{}
	for _, cat := range order {
		v := volume[cat]
		s := sales[cat]
		var vpct, spct float64
		if volumeTotal > 0 {
			vpct = v / volumeTotal * 100
		}
		if salesTotal > 0 {
			spct = s / salesTotal * 100
		}
		res.CategoryByVolume = append(res.CategoryByVolume, CategoryStat{Category: cat, Value: v, Pct: vpct})
		res.CategoryBySales = append(res.CategoryBySales, CategoryStat{Category: cat, Value: s, Pct: spct})
	}
	return nil
}

func (r *ReportRepository) loadTopItemsByCategoryTx(ctx context.Context, tx *sql.Tx, res *ReportSummary, since, until time.Time, where string, args []any) error {
	rows, err := tx.QueryContext(ctx, `
		SELECT cat, name, qty, gross FROM (
			SELECT c.name cat, p.name name, SUM(oi.quantity) qty, SUM(oi.subtotal) gross,
				ROW_NUMBER() OVER (PARTITION BY c.name ORDER BY SUM(oi.quantity) DESC) rn
			FROM orders o
			JOIN order_items oi ON oi.order_id = o.id AND oi.tenant_id=o.tenant_id
			LEFT JOIN products p ON p.id = oi.product_id AND p.tenant_id=o.tenant_id
			LEFT JOIN categories c ON c.id = p.category_id AND c.tenant_id=o.tenant_id
			WHERE `+where+` AND c.id IS NOT NULL
			GROUP BY c.name, p.name
		) t WHERE rn <= 3`, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	res.TopItemsByCategory = []CategoryTopItems{}
	index := map[string]int{}
	for rows.Next() {
		var cat, name string
		var qty int
		var gross float64
		if err := rows.Scan(&cat, &name, &qty, &gross); err != nil {
			return err
		}
		idx, ok := index[cat]
		if !ok {
			idx = len(res.TopItemsByCategory)
			index[cat] = idx
			res.TopItemsByCategory = append(res.TopItemsByCategory, CategoryTopItems{Category: cat, Items: []TopItem{}})
		}
		res.TopItemsByCategory[idx].Items = append(res.TopItemsByCategory[idx].Items, TopItem{Name: name, Qty: qty, GrossSales: gross})
	}
	return rows.Err()
}

// PaymentMethodStat adalah agregasi penjualan per metode pembayaran
// (mendukung split payment karena dihitung dari tabel transactions).
type PaymentMethodStat struct {
	Method       string  `json:"method"`
	Transactions int     `json:"transactions"`
	Total        float64 `json:"total"`
	Pct          float64 `json:"pct"`
}

// PaymentMethods mengembalikan rekap penjualan per metode pembayaran
// pada rentang periode yang sama dengan Summary (plus filter outlet opsional).
func (r *ReportRepository) PaymentMethods(ctx context.Context, period string, outletID *int64) ([]PaymentMethodStat, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)

	since, until, _, _, _, ok := periodRange(period, now)
	if !ok {
		return nil, fmt.Errorf("periode tidak valid: %s", period)
	}

	return withTenantTx1(ctx, r.db, func(tx *sql.Tx, tenantID int64) ([]PaymentMethodStat, error) {
		where := "o.status='completed' AND t.created_at >= $1 AND t.created_at < $2 AND t.tenant_id=$3"
		args := []any{since, until, tenantID}
		if outletID != nil {
			args = append(args, outletUUID(*outletID))
			where += fmt.Sprintf(" AND o.outlet_id = $%d", len(args))
		}

		rows, err := tx.QueryContext(ctx, `
			SELECT COALESCE(NULLIF(t.payment_method,''),'cash') AS method, COUNT(*), COALESCE(SUM(t.amount),0)
			FROM transactions t
			JOIN orders o ON o.id = t.order_id AND o.tenant_id = t.tenant_id
			WHERE `+where+`
			GROUP BY method ORDER BY SUM(t.amount) DESC`, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		list := []PaymentMethodStat{}
		var grandTotal float64
		for rows.Next() {
			var s PaymentMethodStat
			if err := rows.Scan(&s.Method, &s.Transactions, &s.Total); err != nil {
				return nil, err
			}
			grandTotal += s.Total
			list = append(list, s)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if grandTotal > 0 {
			for i := range list {
				list[i].Pct = list[i].Total / grandTotal * 100
			}
		}
		return list, nil
	})
}
