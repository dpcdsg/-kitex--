package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	"github.com/ozline/tiktok/kitex_gen/order"
)

// SeckillAction .
// @router /seckill/action/ [POST]
func SeckillAction(ctx context.Context, c *app.RequestContext) {
	var req api.SeckillActionRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	orderNo, err := rpc.SeckillAction(ctx, &order.SeckillRequest{
		Token:      req.Token,
		ActivityId: req.ActivityID,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.SeckillActionResponse)
	resp.OrderNo = orderNo
	pack.SendResponse(c, resp)
}

// OrderList .
// @router /seckill/order/list/ [GET]
func OrderList(ctx context.Context, c *app.RequestContext) {
	var req api.OrderListRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	page := req.GetPage()
	size := req.GetSize()
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	orders, total, err := rpc.OrderList(ctx, &order.ListOrdersRequest{
		Token:  req.Token,
		Status: req.GetStatus(),
		Page:   page,
		Size:   size,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.OrderListResponse)
	resp.OrderList = pack.OrderList(orders)
	resp.Total = total
	pack.SendResponse(c, resp)
}

// OrderDetail .
// @router /seckill/order/detail/ [GET]
func OrderDetail(ctx context.Context, c *app.RequestContext) {
	var req api.OrderDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	o, err := rpc.OrderDetail(ctx, &order.GetOrderRequest{
		Token:   req.Token,
		OrderNo: req.OrderNo,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.OrderDetailResponse)
	resp.Order = pack.OrderItem(o)
	pack.SendResponse(c, resp)
}

// OrderPay .
// @router /seckill/order/pay/ [POST]
func OrderPay(ctx context.Context, c *app.RequestContext) {
	var req api.OrderPayRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	err := rpc.OrderPay(ctx, &order.PayOrderRequest{
		Token:   req.Token,
		OrderNo: req.OrderNo,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	pack.SendResponse(c, &api.OrderPayResponse{})
}

// OrderCancel .
// @router /seckill/order/cancel/ [POST]
func OrderCancel(ctx context.Context, c *app.RequestContext) {
	var req api.OrderCancelRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	err := rpc.OrderCancel(ctx, &order.CancelOrderRequest{
		Token:   req.Token,
		OrderNo: req.OrderNo,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	pack.SendResponse(c, &api.OrderCancelResponse{})
}
