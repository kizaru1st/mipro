package models

import "time"

type OrderCustomer struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	User       User
	UserID     string `gorm:"size:36;index"`
	Order      Order
	OrderID    string `gorm:"size:36;index"`
	FirstName  string `gorm:"size:100;not null"`
	LastName   string `gorm:"size:100;not null"`
	CityID     string `gorm:"size:100;"`
	ProvinceID string `gorm:"size:100;"`
	Address1   string `gorm:"size:100;"`
	Address2   string `gorm:"size:100;"`
	Phone      string `gorm:"size:50;"`
	Email      string `gorm:"size:100;"`
	PostCode   string `gorm:"size:100;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderCustomerResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	PostCode   string `json:"post_code"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
