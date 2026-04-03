package api

import (
	"context"
	"strconv"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/middleware/metrics"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	orderdb "github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/kitex_gen/order"
	"github.com/ozline/tiktok/pkg/errno"
	"github.com/ozline/tiktok/pkg/utils"
)

type seckillActionRequest struct {
	Token     string `json:"token" form:"token" query:"token"`
	ActivityID string `json:"activity_id" form:"activity_id" query:"activity_id"`
}

// SeckillAction .
// @router /seckill/action/ [POST]
func SeckillAction(ctx context.Context, c *app.RequestContext) {
	var bizErr error
	defer func() { metrics.RecordBusinessResult(metrics.HandlerSeckillAction, bizErr) }()

	var req seckillActionRequest
	if err := c.BindAndValidate(&req); err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	activityID, err := strconv.ParseInt(req.ActivityID, 10, 64)
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	orderNo, err := rpc.SeckillAction(ctx, &order.SeckillRequest{
		Token:      req.Token,
		ActivityId: activityID,
	})
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.SeckillActionResponse)
	resp.OrderNo = orderNo
	pack.SendResponse(c, resp)
}

type orderBuyRequest struct {
	Token     string `json:"token" form:"token" query:"token"`
	ProductID string `json:"product_id" form:"product_id" query:"product_id"`
}

var initOrderDBOnce sync.Once

// OrderBuy .
// @router /seckill/order/buy/ [POST]
func OrderBuy(ctx context.Context, c *app.RequestContext) {
	var bizErr error
	defer func() { metrics.RecordBusinessResult(metrics.HandlerOrderBuy, bizErr) }()

	var req orderBuyRequest
	if err := c.BindAndValidate(&req); err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	// 1) 查询商品信息（包含价格/库存）
	p, err := rpc.ProductDetail(ctx, productID, req.Token)
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}
	if p == nil || p.Stock <= 0 {
		bizErr = errno.StockSoldOutError
		pack.SendFailResponse(c, errno.StockSoldOutError)
		return
	}

	// 2) 扣减普通商品库存（复用 ProductUpdate）
	newStock := p.Stock - 1
	_, err = rpc.ProductUpdate(ctx, &product.UpdateProductRequest{
		Token:       req.Token,
		ProductId:   productID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       newStock,
		ImageUrl:    p.ImageUrl,
		Category:    p.Category,
	})
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	// 3) 创建订单
	initOrderDBOnce.Do(func() { orderdb.Init() })
	order, err := orderdb.CreateOrder(ctx, &orderdb.Order{
		UserId:     claims.UserId,
		ProductId:  productID,
		ActivityId: 0,
		Amount:     p.Price,
	})
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	pack.SendResponse(c, map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"order_no":    order.OrderNo,
	})
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
