package main

import (
	"context"

	"github.com/ozline/tiktok/cmd/order/pack"
	"github.com/ozline/tiktok/cmd/order/service"
	"github.com/ozline/tiktok/kitex_gen/order"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/ozline/tiktok/pkg/utils"
)

type OrderServiceImpl struct{}

func (s *OrderServiceImpl) Seckill(ctx context.Context, req *order.SeckillRequest) (resp *order.SeckillResponse, err error) {
	resp = new(order.SeckillResponse)

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, nil
	}

	if req.ActivityId <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	orderNo, err := service.NewOrderService(ctx).Seckill(claims.UserId, req.ActivityId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.OrderNo = orderNo
	return
}

func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderRequest) (resp *order.GetOrderResponse, err error) {
	resp = new(order.GetOrderResponse)

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, nil
	}

	if len(req.OrderNo) == 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	o, err := service.NewOrderService(ctx).GetOrder(claims.UserId, req.OrderNo)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Order = pack.Order(o)
	return
}

func (s *OrderServiceImpl) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (resp *order.ListOrdersResponse, err error) {
	resp = new(order.ListOrdersResponse)

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, nil
	}

	orders, total, err := service.NewOrderService(ctx).ListOrders(claims.UserId, req.Status, req.Page, req.Size)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.OrderList = pack.Orders(orders)
	resp.Total = total
	return
}

func (s *OrderServiceImpl) PayOrder(ctx context.Context, req *order.PayOrderRequest) (resp *order.PayOrderResponse, err error) {
	resp = new(order.PayOrderResponse)

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, nil
	}

	if len(req.OrderNo) == 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	err = service.NewOrderService(ctx).PayOrder(claims.UserId, req.OrderNo)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	return
}

func (s *OrderServiceImpl) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (resp *order.CancelOrderResponse, err error) {
	resp = new(order.CancelOrderResponse)

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = pack.BuildBaseResp(errno.AuthorizationFailedError)
		return resp, nil
	}

	if len(req.OrderNo) == 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	err = service.NewOrderService(ctx).CancelOrder(claims.UserId, req.OrderNo)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	return
}
