package service

import (
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/pkg/errno"
)

func (s *OrderService) GetOrder(userId int64, orderNo string) (*db.Order, error) {
	order, err := db.GetOrderByNo(s.ctx, orderNo)
	if err != nil {
		return nil, err
	}
	if order.UserId != userId {
		return nil, errno.OrderNotFoundError
	}
	return order, nil
}
