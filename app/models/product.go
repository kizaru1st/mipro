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

type ProductResponse struct {
	ID               string          `json:"id"`
	ParentID         string          `json:"parent_id"`
	UserID           string          `json:"user_id"`
	Sku              string          `json:"sku"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Price            decimal.Decimal `json:"price"`
	Stock            int             `json:"stock"`
	Weight           decimal.Decimal `json:"weight"`
	ShortDescription string          `json:"short_description"`
	Description      string          `json:"description"`
	Status           int             `json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `json:"deleted_at"`
	ProductImage     []ProductImage  `json:"product_image"`
	Categories       []Category      `json:"categories"`
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

func (p *Product) FindByID(db *gorm.DB, productID string) (*ProductResponse, error) {
	var err error
	var product Product

	err = db.Debug().Model(&Product{}).Where("id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}

	// Buat instansi ProductResponse dan isi dengan data dari product
	productResponse := ProductResponse{
		ID:               product.ID,
		ParentID:         product.ParentID,
		UserID:           product.UserID,
		Sku:              product.Sku,
		Name:             product.Name,
		Slug:             product.Slug,
		Price:            product.Price,
		Stock:            product.Stock,
		Weight:           product.Weight,
		ShortDescription: product.ShortDescription,
		Description:      product.Description,
		Status:           product.Status,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
		DeletedAt:        product.DeletedAt,
		ProductImage:     product.ProductImage,
		Categories:       product.Categories,
	}

	return &productResponse, nil
}
