package models

import "time"

type ProductImage struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Product    Product
	ProductID  string `gorm:"size:36;index"`
	Path       string `gorm:"type:text"`
	ExtraLarge string `gorm:"type:text"`
	Large      string `gorm:"type:text"`
	Medium     string `gorm:"type:text"`
	Small      string `gorm:"type:text"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProductImageResponse struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	Path       string `json:"path"`
	ExtraLarge string `json:"extra_large"`
	Large      string `json:"large"`
	Medium     string `json:"medium"`
	Small      string `json:"small"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
