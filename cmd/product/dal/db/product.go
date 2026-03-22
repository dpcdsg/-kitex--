package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id          int64          `gorm:"primaryKey"`
	SellerId    int64          `gorm:"not null;index;default:0"`
	Name        string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Price       int64          `gorm:"not null"`
	Stock       int64          `gorm:"not null;default:0"`
	ImageUrl    string         `gorm:"type:varchar(512)"`
	Category    string         `gorm:"type:varchar(128);index"`
	Status      int64          `gorm:"default:1"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "product"
}

func CreateProduct(ctx context.Context, p *Product) (*Product, error) {
	p.Id = SF.NextVal()
	err := DB.WithContext(ctx).Create(p).Error
	return p, err
}

func GetProductByID(ctx context.Context, id int64) (*Product, error) {
	p := &Product{}
	err := DB.WithContext(ctx).Where("id = ?", id).First(p).Error
	return p, err
}

func ListProducts(ctx context.Context, page, size int64, category string) ([]*Product, int64, error) {
	var products []*Product
	var total int64

	query := DB.WithContext(ctx).Model(&Product{})
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(int(offset)).Limit(int(size)).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func UpdateProduct(ctx context.Context, p *Product) error {
	return DB.WithContext(ctx).Model(&Product{}).Where("id = ?", p.Id).Updates(map[string]interface{}{
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"image_url":   p.ImageUrl,
		"category":    p.Category,
	}).Error
}
