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
	// activity.StartTime/EndTime 返回的是“不带时区”的字符串。
	// 若直接 time.Parse(...)，Go 会把它当作 UTC，可能导致与服务端本地时区比较时提前/延后判定。
	// 使用 ParseInLocation 按服务端本地时区解释这些时间字符串。
	const layout = "2006-01-02 15:04:05"
	startTime, _ := time.ParseInLocation(layout, activity.StartTime, time.Local)
	endTime, _ := time.ParseInLocation(layout, activity.EndTime, time.Local)

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
