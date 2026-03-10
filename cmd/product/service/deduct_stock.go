package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/cache"
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/pkg/errno"
)

func (s *ProductService) DeductStock(activityId int64) error {
	ok, err := cache.DeductSeckillStock(s.ctx, activityId)
	if err != nil {
		return errno.ServiceInternalError.WithMessage("redis deduct stock failed")
	}
	if !ok {
		return errno.NewErrNo(errno.ParamErrorCode, "stock sold out")
	}

	if err := db.DeductStock(s.ctx, activityId); err != nil {
		return errno.ServiceInternalError.WithMessage("db deduct stock failed")
	}

	return nil
}
