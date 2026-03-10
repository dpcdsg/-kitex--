package pack

import (
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/kitex_gen/order"
)

func OrderItem(o *order.Order) *api.Order {
	if o == nil {
		return nil
	}
	return &api.Order{
		ID:          o.Id,
		OrderNo:     o.OrderNo,
		UserID:      o.UserId,
		ProductID:   o.ProductId,
		ProductName: o.ProductName,
		ActivityID:  o.ActivityId,
		Amount:      o.Amount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
	}
}

func OrderList(orders []*order.Order) []*api.Order {
	result := make([]*api.Order, 0, len(orders))
	for _, o := range orders {
		result = append(result, OrderItem(o))
	}
	return result
}
