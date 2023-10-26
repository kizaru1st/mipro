package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Shipment struct {
	ID          string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	User        User
	UserID      string `gorm:"size:36;index"`
	Order       Order
	OrderID     string `gorm:"size:36;index"`
	TrackNumber string `gorm:"size:255;index"`
	Status      string `gorm:"size:36;index"`
	TotalQty    int
	TotalWeight decimal.Decimal `gorm:"type:decimal(10,2);"`
	FirstName   string          `gorm:"size:100;not null"`
	LastName    string          `gorm:"size:100;not null"`
	CityID      string          `gorm:"size:100;"`
	ProvinceID  string          `gorm:"size:100;"`
	Address1    string          `gorm:"size:100;"`
	Address2    string          `gorm:"size:100;"`
	Phone       string          `gorm:"size:50;"`
	Email       string          `gorm:"size:100;"`
	PostCode    string          `gorm:"size:100;"`
	ShippedBy   string          `gorm:"size:36"`
	ShippedAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type ShipmentResponse struct {
	ID          string `json:"id"`
	TrackNumber string `json:"track_number"`
	Status      string `json:"status"`
	TotalQty    int    `json:"total_qty"`
	TotalWeight string `json:"total_weight"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CityID      string `json:"city_id"`
	ProvinceID  string `json:"province_id"`
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	PostCode    string `json:"post_code"`
	ShippedBy   string `json:"shipped_by"`
	ShippedAt   string `json:"shipped_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
