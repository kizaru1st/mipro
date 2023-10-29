package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shopspring/decimal"
)

type CartItem struct {
	ID              string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Cart            Cart
	CartID          string `gorm:"size:36;index"`
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
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CartItemResponse struct {
	ID              string          `json:"id"`
	CartID          string          `json:"cart_id"`
	ProductID       string          `json:"product_id"`
	Qty             int             `json:"qty"`
	BasePrice       decimal.Decimal `json:"base_price"`
	BaseTotal       decimal.Decimal `json:"base_total"`
	TaxAmount       decimal.Decimal `json:"tax_amount"`
	TaxPercent      decimal.Decimal `json:"tax_percent"`
	DiscountAmount  decimal.Decimal `json:"discount_amount"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
	SubTotal        decimal.Decimal `json:"sub_total"`
}

func (c *CartItem) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	return nil
}
