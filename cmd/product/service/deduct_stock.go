package service

import (
	"errors"

	"github.com/ozline/tiktok/cmd/product/dal/cache"
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func (s *ProductService) DeductStock(activityId int64) error {
	ok, productId, err := cache.DeductSeckillStock(s.ctx, activityId)
	if err != nil {
		// 兼容旧活动：如果 redis 里缺少 seckill->product 映射，则补齐 key 再重试一次
		if errors.Is(err, redis.Nil) {
			activity, dbErr := db.GetSeckillByID(s.ctx, activityId)
			if dbErr != nil {
				return errno.ParamError.WithMessage("seckill activity not found")
			}
			prod, dbErr := db.GetProductByID(s.ctx, activity.ProductId)
			if dbErr != nil {
				return errno.ParamError.WithMessage("product not found")
			}
			_ = cache.WarmUpProductStock(s.ctx, prod.Id, prod.Stock)
			_ = cache.WarmUpSeckillStock(s.ctx, activity.Id, activity.ProductId, activity.AvailableStock)

			ok, productId, err = cache.DeductSeckillStock(s.ctx, activityId)
		}
	}

	if err != nil {
		return errno.ServiceInternalError.WithMessage("redis deduct stock failed")
	}
	if !ok {
		return errno.NewErrNo(errno.ParamErrorCode, "stock sold out")
	}

	if err := db.DeductStock(s.ctx, activityId, productId); err != nil {
		// 若 redis 侧成功但 DB 侧失败（可能存在并发/漂移），需要回滚 redis 计数，避免“redis 锁住库存但 DB 没扣”。
		_ = cache.RevertSeckillAndProductStock(s.ctx, activityId, productId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.StockSoldOutError
		}
		return errno.ServiceInternalError.WithMessage("db deduct stock failed")
	}

	return nil
}
