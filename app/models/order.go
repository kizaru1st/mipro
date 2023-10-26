package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	ID                  string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	UserID              string `gorm:"size:36;index"`
	User                User
	OrderItems          []OrderItem
	OrderCustomer       *OrderCustomer
	Code                string `gorm:"size:50;index"`
	Status              int
	OrderDate           time.Time
	PaymentDue          time.Time
	PaymentStatus       string          `gorm:"size:50;index"`
	PaymentToken        string          `gorm:"size:100;index"`
	BaseTotalPrice      decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxAmount           decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxPercent          decimal.Decimal `gorm:"type:decimal(10,2)"`
	DiscountAmount      decimal.Decimal `gorm:"type:decimal(16,2)"`
	DiscountPercent     decimal.Decimal `gorm:"type:decimal(10,2)"`
	ShippingCost        decimal.Decimal `gorm:"type:decimal(16,2)"`
	GrandTotal          decimal.Decimal `gorm:"type:decimal(16,2)"`
	Note                string          `gorm:"type:text"`
	ShippingCourier     string          `gorm:"size:100"`
	ShippingServiceName string          `gorm:"size:100"`
	ApprovedBy          string          `gorm:"size:36"`
	ApprovedAt          time.Time
	CancelledBy         string `gorm:"size:36"`
	CancelledAt         time.Time
	CancellationNote    string `gorm:"size:255"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

type OrderResponse struct {
	ID                  string `json:"id"`
	UserID              string `json:"user_id"`
	Code                string `json:"code"`
	Status              int    `json:"status"`
	OrderDate           string `json:"order_date"`
	PaymentDue          string `json:"payment_due"`
	PaymentStatus       string `json:"payment_status"`
	PaymentToken        string `json:"payment_token"`
	BaseTotalPrice      string `json:"base_total_price"`
	TaxAmount           string `json:"tax_amount"`
	TaxPercent          string `json:"tax_percent"`
	DiscountAmount      string `json:"discount_amount"`
	DiscountPercent     string `json:"discount_percent"`
	ShippingCost        string `json:"shipping_cost"`
	GrandTotal          string `json:"grand_total"`
	Note                string `json:"note"`
	ShippingCourier     string `json:"shipping_courier"`
	ShippingServiceName string `json:"shipping_service_name"`
	ApprovedBy          string `json:"approved_by"`
	ApprovedAt          string `json:"approved_at"`
	CancelledBy         string `json:"cancelled_by"`
	CancelledAt         string `json:"cancelled_at"`
	CancellationNote    string `json:"cancellation_note"`
}
