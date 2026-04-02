package service

import (
	"time"

	"github.com/ozline/tiktok/cmd/product/dal/cache"
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/ozline/tiktok/pkg/utils"
)

func (s *ProductService) CreateSeckill(req *product.CreateSeckillRequest) (*db.SeckillActivity, error) {
	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		return nil, errno.AuthorizationFailedError
	}

	prod, err := db.GetProductByID(s.ctx, req.ProductId)
	if err != nil {
		return nil, errno.ParamError.WithMessage("product not found")
	}
	if prod.SellerId != claims.UserId {
		return nil, errno.ProductPermissionDeniedError
	}

	// 秒杀总库存以普通商品库存为上限，避免“商品库存不够但秒杀仍可卖出”的问题。
	if req.TotalStock > prod.Stock {
		return nil, errno.ParamError.WithMessage("秒杀总库存不能大于商品库存")
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errno.ParamError.WithMessage("invalid start_time format")
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, errno.ParamError.WithMessage("invalid end_time format")
	}

	activity := &db.SeckillActivity{
		ProductId:      req.ProductId,
		SeckillPrice:   req.SeckillPrice,
		TotalStock:     req.TotalStock,
		AvailableStock: req.TotalStock,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         0,
	}

	result, err := db.CreateSeckill(s.ctx, activity)
	if err != nil {
		return nil, err
	}

	_ = cache.WarmUpProductStock(s.ctx, prod.Id, prod.Stock)
	_ = cache.WarmUpSeckillStock(s.ctx, result.Id, result.ProductId, result.AvailableStock)

	return result, nil
}
