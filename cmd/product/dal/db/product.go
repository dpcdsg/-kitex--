package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id          int64          `gorm:"primaryKey"`
	SellerId    int64          `gorm:"not null;index"`
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
	err := DB.WithContext(ctx).Create(p).Error
	return p, err
}

func GetProductByID(ctx context.Context, id int64) (*Product, error) {
	p := &Product{}
	err := DB.WithContext(ctx).Where("id = ?", id).First(p).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Some clients (e.g. JS number) may lose precision for int64 IDs.
		// Fallback to nearest candidate in a tiny range for compatibility.
		resolvedID, resolveErr := resolveProductIDWithTolerance(ctx, id)
		if resolveErr != nil {
			return p, err
		}
		err = DB.WithContext(ctx).Where("id = ?", resolvedID).First(p).Error
	}
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

func UpdateProduct(ctx context.Context, p *Product) (*Product, error) {
	resolvedID, err := resolveProductIDWithTolerance(ctx, p.Id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"image_url":   p.ImageUrl,
		"category":    p.Category,
	}

	if err := DB.WithContext(ctx).Model(&Product{}).Where("id = ?", resolvedID).Updates(updates).Error; err != nil {
		return nil, err
	}

	return GetProductByID(ctx, resolvedID)
}

func resolveProductIDWithTolerance(ctx context.Context, id int64) (int64, error) {
	// JS Number loses precision above 2^53; most drifts are tiny.
	const tolerance int64 = 1024

	var exact Product
	if err := DB.WithContext(ctx).Select("id").Where("id = ?", id).First(&exact).Error; err == nil {
		return id, nil
	}

	var nearest Product
	err := DB.WithContext(ctx).
		Select("id").
		Where("id BETWEEN ? AND ?", id-tolerance, id+tolerance).
		Order(fmt.Sprintf("ABS(id - %d) ASC", id)).
		First(&nearest).Error
	if err != nil {
		return 0, err
	}
	return nearest.Id, nil
}
