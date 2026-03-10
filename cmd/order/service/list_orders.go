package service

import (
	"github.com/ozline/tiktok/cmd/order/dal/db"
)

func (s *OrderService) ListOrders(userId, status, page, size int64) ([]*db.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return db.ListOrders(s.ctx, userId, status, page, size)
}
