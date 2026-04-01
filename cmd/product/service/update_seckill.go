package service

import (
	"time"

	"github.com/ozline/tiktok/cmd/product/dal/cache"
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/errno"
)

func (s *ProductService) UpdateSeckill(req *product.UpdateSeckillRequest) (*product.SeckillActivity, error) {
	activity, err := db.GetSeckillByID(s.ctx, req.ActivityId)
	if err != nil {
		return nil, errno.ParamError.WithMessage("seckill activity not found")
	}

	p, err := db.GetProductByID(s.ctx, req.ProductId)
	if err != nil {
		return nil, errno.ParamError.WithMessage("product not found")
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errno.ParamError.WithMessage("invalid start_time format")
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, errno.ParamError.WithMessage("invalid end_time format")
	}
	if !endTime.After(startTime) {
		return nil, errno.ParamError.WithMessage("end_time must be after start_time")
	}

	// Keep sold quantity consistent when total stock changes.
	sold := activity.TotalStock - activity.AvailableStock
	if sold < 0 {
		sold = 0
	}
	if req.TotalStock < sold {
		return nil, errno.ParamError.WithMessage("total_stock is less than sold quantity")
	}
	newAvailable := req.TotalStock - sold

	now := time.Now()
	status := int64(1)
	if now.Before(startTime) {
		status = 0
	} else if now.After(endTime) || newAvailable <= 0 {
		status = 2
	}

	updated, err := db.UpdateSeckill(s.ctx, req.ActivityId, map[string]interface{}{
		"product_id":       req.ProductId,
		"seckill_price":    req.SeckillPrice,
		"total_stock":      req.TotalStock,
		"available_stock":  newAvailable,
		"start_time":       startTime,
		"end_time":         endTime,
		"status":           status,
	})
	if err != nil {
		return nil, err
	}

	_ = cache.WarmUpSeckillStock(s.ctx, updated.Id, updated.AvailableStock)

	return &product.SeckillActivity{
		Id:             updated.Id,
		ProductId:      updated.ProductId,
		ProductName:    p.Name,
		SeckillPrice:   updated.SeckillPrice,
		TotalStock:     updated.TotalStock,
		AvailableStock: updated.AvailableStock,
		StartTime:      updated.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:        updated.EndTime.Format("2006-01-02 15:04:05"),
		Status:         updated.Status,
	}, nil
}
