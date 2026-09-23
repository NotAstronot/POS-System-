package entity

import "time"

type Product struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	BasePrice     float64   `json:"base_price"`
	Unit          string    `json:"unit"`
	ImageURL      string    `json:"image_url"`
	CategoryID    *int64    `json:"category_id"`
	HasVariants   bool      `json:"has_variants"`
	TaxRate       float64   `json:"tax_rate"`
	SortOrder     int       `json:"sort_order"`
	Stock         int       `json:"stock"`
	Barcode       string    `json:"barcode"`
	PurchasePrice float64   `json:"purchase_price"`
	MinStock      int       `json:"min_stock"`
	IsService     bool      `json:"is_service"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductAvailability struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	Date        string    `json:"date"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}
