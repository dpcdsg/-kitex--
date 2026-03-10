package service

import (
	"fmt"
	"time"

	"github.com/ozline/tiktok/cmd/order/dal/cache"
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/cmd/order/dal/mq"
	"github.com/ozline/tiktok/cmd/order/rpc"
	"github.com/ozline/tiktok/pkg/errno"
)

// Seckill 秒杀下单核心流程：
// 1. Redis防重复下单（SetNX）
// 2. RPC查询活动状态
// 3. RPC调用Redis+Lua原子扣库存
// 4. 创建订单记录
// 5. MQ异步通知（削峰填谷）
func (s *OrderService) Seckill(userId, activityId int64) (string, error) {
	ok, err := cache.CheckSeckillRepeat(s.ctx, userId, activityId)
	if err != nil {
		return "", errno.ServiceInternalError.WithMessage("check repeat failed")
	}
	if !ok {
		return "", errno.RepeatSeckillError
	}

	activity, err := rpc.GetSeckillActivity(s.ctx, activityId, "")
	if err != nil {
		cache.RemoveSeckillRepeatFlag(s.ctx, userId, activityId)
		return "", err
	}

	now := time.Now()
	startTime, _ := time.Parse("2006-01-02 15:04:05", activity.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", activity.EndTime)

	if now.Before(startTime) {
		cache.RemoveSeckillRepeatFlag(s.ctx, userId, activityId)
		return "", errno.SeckillNotStartedError
	}
	if now.After(endTime) {
		cache.RemoveSeckillRepeatFlag(s.ctx, userId, activityId)
		return "", errno.SeckillEndedError
	}

	err = rpc.DeductStock(s.ctx, activityId, "")
	if err != nil {
		cache.RemoveSeckillRepeatFlag(s.ctx, userId, activityId)
		return "", err
	}

	order, err := db.CreateOrder(s.ctx, &db.Order{
		UserId:     userId,
		ProductId:  activity.ProductId,
		ActivityId: activityId,
		Amount:     activity.SeckillPrice,
	})
	if err != nil {
		cache.RemoveSeckillRepeatFlag(s.ctx, userId, activityId)
		return "", errno.ServiceInternalError.WithMessage(fmt.Sprintf("create order failed: %v", err))
	}

	_ = mq.PublishSeckillOrder(s.ctx, &mq.SeckillMessage{
		UserId:     userId,
		ActivityId: activityId,
		OrderNo:    order.OrderNo,
	})

	return order.OrderNo, nil
}
