package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID               string `gorm:"size:100;not null;uniqueIndex;primary_key"`
	ParentID         string `gorm:"size:100;index"`
	User             User
	UserID           string          `gorm:"size:100;index"`
	Sku              string          `gorm:"size:100;index"`
	Name             string          `gorm:"size:100"`
	Slug             string          `gorm:"size:100"`
	Price            decimal.Decimal `gorm:"type:decimal(16,2)"`
	Stock            int
	Weight           decimal.Decimal `gorm:"type:decimal(10,2)"`
	ShortDescription string          `gorm:"size:100"`
	Description      string          `gorm:"type:text"`
	Status           int             `gorm:"default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
	ProductImage     []ProductImage
	Categories       []Category `gorm:"many2many:product_categories;"`
}

func (p *Product) GetProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	var products []Product

	err = db.Debug().Limit(20).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &products, nil
}

func (p *Product) FindBySlug(db *gorm.DB, slug string) (*Product, error) {
	var err error
	var product Product

	err = db.Debug().Model(&Product{}).Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}
