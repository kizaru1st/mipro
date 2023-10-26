package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID              string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Order           Order
	OrderID         string `gorm:"size:36;index"`
	Product         Product
	ProductID       string `gorm:"size:36;index"`
	Qty             int
	BasePrice       decimal.Decimal `gorm:"type:decimal(16,2)"`
	BaseTotal       decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxPercent      decimal.Decimal `gorm:"type:decimal(10,2)"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(16,2)"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(10,2)"`
	SubTotal        decimal.Decimal `gorm:"type:decimal(16,2)"`
	Sku             string          `gorm:"size:36;index"`
	Name            string          `gorm:"size:255"`
	Weight          decimal.Decimal `gorm:"type:decimal(10,2)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItemResponse struct {
	ID              string `json:"id"`
	OrderID         string `json:"order_id"`
	ProductID       string `json:"product_id"`
	Qty             int    `json:"qty"`
	BasePrice       string `json:"base_price"`
	BaseTotal       string `json:"base_total"`
	TaxAmount       string `json:"tax_amount"`
	TaxPercent      string `json:"tax_percent"`
	DiscountAmount  string `json:"discount_amount"`
	DiscountPercent string `json:"discount_percent"`
	SubTotal        string `json:"sub_total"`
	Sku             string `json:"sku"`
	Name            string `json:"name"`
	Weight          string `json:"weight"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
