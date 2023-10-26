package models

import "time"

type Section struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Name       string `gorm:"size:100;"`
	Slug       string `gorm:"size:100;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Categories []Category
}

type SectionResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
