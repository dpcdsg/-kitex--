package service

import "context"

type OrderService struct {
	ctx context.Context
}

func NewOrderService(ctx context.Context) *OrderService {
	return &OrderService{ctx: ctx}
}
