package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	"github.com/ozline/tiktok/kitex_gen/product"
)

// ProductCreate .
// @router /seckill/product/create/ [POST]
func ProductCreate(ctx context.Context, c *app.RequestContext) {
	var req api.ProductCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	p, err := rpc.ProductCreate(ctx, &product.CreateProductRequest{
		Token:       req.Token,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageUrl,
		Category:    req.Category,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.ProductCreateResponse)
	resp.Product = pack.Product(p)
	pack.SendResponse(c, resp)
}

// ProductDetail .
// @router /seckill/product/detail/ [GET]
func ProductDetail(ctx context.Context, c *app.RequestContext) {
	var req api.ProductDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	p, err := rpc.ProductDetail(ctx, req.ProductID, req.GetToken())
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.ProductDetailResponse)
	resp.Product = pack.Product(p)
	pack.SendResponse(c, resp)
}

// ProductList .
// @router /seckill/product/list/ [GET]
func ProductList(ctx context.Context, c *app.RequestContext) {
	var req api.ProductListRequest
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

	products, total, err := rpc.ProductList(ctx, &product.ListProductsRequest{
		Token:    req.GetToken(),
		Page:     page,
		Size:     size,
		Category: req.GetCategory(),
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.ProductListResponse)
	resp.ProductList = pack.ProductList(products)
	resp.Total = total
	pack.SendResponse(c, resp)
}

// SeckillCreate .
// @router /seckill/activity/create/ [POST]
func SeckillCreate(ctx context.Context, c *app.RequestContext) {
	var req api.SeckillCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activity, err := rpc.SeckillCreate(ctx, &product.CreateSeckillRequest{
		Token:        req.Token,
		ProductId:    req.ProductID,
		SeckillPrice: req.SeckillPrice,
		TotalStock:   req.TotalStock,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.SeckillCreateResponse)
	resp.Activity = pack.SeckillActivity(activity)
	pack.SendResponse(c, resp)
}

// SeckillList .
// @router /seckill/activity/list/ [GET]
func SeckillList(ctx context.Context, c *app.RequestContext) {
	var req api.SeckillListRequest
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

	activities, total, err := rpc.SeckillList(ctx, &product.ListSeckillRequest{
		Token:  req.GetToken(),
		Status: req.GetStatus(),
		Page:   page,
		Size:   size,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.SeckillListResponse)
	resp.ActivityList = pack.SeckillActivityList(activities)
	resp.Total = total
	pack.SendResponse(c, resp)
}

// SeckillDetail .
// @router /seckill/activity/detail/ [GET]
func SeckillDetail(ctx context.Context, c *app.RequestContext) {
	var req api.SeckillDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activity, err := rpc.SeckillDetail(ctx, req.ActivityID, req.GetToken())
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.SeckillDetailResponse)
	resp.Activity = pack.SeckillActivity(activity)
	pack.SendResponse(c, resp)
}
