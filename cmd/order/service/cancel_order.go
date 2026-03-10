package service

import (
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/pkg/constants"
	"github.com/ozline/tiktok/pkg/errno"
)

func (s *OrderService) CancelOrder(userId int64, orderNo string) error {
	order, err := db.GetOrderByNo(s.ctx, orderNo)
	if err != nil {
		return err
	}
	if order.UserId != userId {
		return errno.OrderNotFoundError
	}
	if order.Status != constants.OrderStatusPending {
		return errno.OrderStatusError
	}

	return db.UpdateOrderStatus(s.ctx, orderNo, constants.OrderStatusCancelled)
}
