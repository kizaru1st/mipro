package models

import "time"

type Category struct {
	ID        string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	ParentID  string `gorm:"size:36;"`
	Section   Section
	SectionID string    `gorm:"size:36;index"`
	Products  []Product `gorm:"many2many:product_categories;"`
	Name      string    `gorm:"size:100;"`
	Slug      string    `gorm:"size:100;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CategoryResponse struct {
	ID        string `json:"id"`
	ParentID  string `json:"parent_id"`
	SectionID string `json:"section_id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
}
