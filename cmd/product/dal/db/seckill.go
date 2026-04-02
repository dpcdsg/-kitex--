package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SeckillActivity struct {
	Id             int64          `gorm:"primaryKey"`
	ProductId      int64          `gorm:"not null;index"`
	SeckillPrice   int64          `gorm:"not null"`
	TotalStock     int64          `gorm:"not null"`
	AvailableStock int64          `gorm:"not null"`
	StartTime      time.Time      `gorm:"not null"`
	EndTime        time.Time      `gorm:"not null"`
	Status         int64          `gorm:"default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (SeckillActivity) TableName() string {
	return "seckill_activity"
}

func CreateSeckill(ctx context.Context, s *SeckillActivity) (*SeckillActivity, error) {
	s.Id = SF.NextVal()
	err := DB.WithContext(ctx).Create(s).Error
	return s, err
}

func GetSeckillByID(ctx context.Context, id int64) (*SeckillActivity, error) {
	s := &SeckillActivity{}
	err := DB.WithContext(ctx).Where("id = ?", id).First(s).Error
	return s, err
}

func ListSeckills(ctx context.Context, status, page, size int64) ([]*SeckillActivity, int64, error) {
	var activities []*SeckillActivity
	var total int64

	query := DB.WithContext(ctx).Model(&SeckillActivity{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Offset(int(offset)).Limit(int(size)).Order("created_at DESC").Find(&activities).Error; err != nil {
		return nil, 0, err
	}

	return activities, total, nil
}

func DeductStock(ctx context.Context, activityId int64, productId int64) error {
	// 关键：秒杀扣减成功后，同步扣减商品真实库存，防止“秒杀超卖但商品库存没变”。
	tx := DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	res1 := tx.Model(&SeckillActivity{}).
		Where("id = ? AND available_stock > 0", activityId).
		Update("available_stock", gorm.Expr("available_stock - 1"))
	if res1.Error != nil {
		_ = tx.Rollback()
		return res1.Error
	}
	if res1.RowsAffected == 0 {
		_ = tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	res2 := tx.Model(&Product{}).
		Where("id = ? AND stock > 0", productId).
		Update("stock", gorm.Expr("stock - 1"))
	if res2.Error != nil {
		_ = tx.Rollback()
		return res2.Error
	}
	if res2.RowsAffected == 0 {
		_ = tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	return tx.Commit().Error
}

func UpdateSeckill(ctx context.Context, id int64, updates map[string]interface{}) (*SeckillActivity, error) {
	if err := DB.WithContext(ctx).Model(&SeckillActivity{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return GetSeckillByID(ctx, id)
}
