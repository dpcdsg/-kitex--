package api

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	"github.com/ozline/tiktok/kitex_gen/product"
)

type productDetailRequest struct {
	Token     string `json:"token" form:"token" query:"token"`
	ProductID string `json:"product_id" form:"product_id" query:"product_id"`
}

type productUpdateRequest struct {
	Token       string `json:"token" form:"token" query:"token"`
	ProductID   string `json:"product_id" form:"product_id" query:"product_id"`
	Name        string `json:"name" form:"name" query:"name"`
	Description string `json:"description" form:"description" query:"description"`
	Price       int64  `json:"price" form:"price" query:"price"`
	Stock       int64  `json:"stock" form:"stock" query:"stock"`
	ImageURL    string `json:"image_url" form:"image_url" query:"image_url"`
	Category    string `json:"category" form:"category" query:"category"`
}

type seckillUpdateRequest struct {
	Token        string `json:"token" form:"token" query:"token"`
	ActivityID   string `json:"activity_id" form:"activity_id" query:"activity_id"`
	ProductID    string `json:"product_id" form:"product_id" query:"product_id"`
	SeckillPrice int64  `json:"seckill_price" form:"seckill_price" query:"seckill_price"`
	TotalStock   int64  `json:"total_stock" form:"total_stock" query:"total_stock"`
	StartTime    string `json:"start_time" form:"start_time" query:"start_time"`
	EndTime      string `json:"end_time" form:"end_time" query:"end_time"`
}

type seckillCreateRequest struct {
	Token        string `json:"token" form:"token" query:"token"`
	ProductID    string `json:"product_id" form:"product_id" query:"product_id"`
	SeckillPrice int64  `json:"seckill_price" form:"seckill_price" query:"seckill_price"`
	TotalStock   int64  `json:"total_stock" form:"total_stock" query:"total_stock"`
	StartTime    string `json:"start_time" form:"start_time" query:"start_time"`
	EndTime      string `json:"end_time" form:"end_time" query:"end_time"`
}

type seckillDetailRequest struct {
	Token      string `json:"token" form:"token" query:"token"`
	ActivityID string `json:"activity_id" form:"activity_id" query:"activity_id"`
}

func productToJSON(p *product.Product) map[string]interface{} {
	if p == nil {
		return nil
	}
	return map[string]interface{}{
		"id":          strconv.FormatInt(p.Id, 10),
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"image_url":   p.ImageUrl,
		"category":    p.Category,
		"seller_id":   strconv.FormatInt(p.SellerId, 10),
	}
}

func productsToJSON(products []*product.Product) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		result = append(result, productToJSON(p))
	}
	return result
}

func seckillToJSON(a *product.SeckillActivity) map[string]interface{} {
	if a == nil {
		return nil
	}
	return map[string]interface{}{
		"id":              strconv.FormatInt(a.Id, 10),
		"product_id":      strconv.FormatInt(a.ProductId, 10),
		"product_name":    a.ProductName,
		"seckill_price":   a.SeckillPrice,
		"total_stock":     a.TotalStock,
		"available_stock": a.AvailableStock,
		"start_time":      a.StartTime,
		"end_time":        a.EndTime,
		"status":          a.Status,
	}
}

func seckillListToJSON(items []*product.SeckillActivity) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		result = append(result, seckillToJSON(it))
	}
	return result
}

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
		ImageUrl:    req.GetImageURL(),
		Category:    req.Category,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"product":     productToJSON(p),
	}
	pack.SendResponse(c, resp)
}

// ProductDetail .
// @router /seckill/product/detail/ [GET]
func ProductDetail(ctx context.Context, c *app.RequestContext) {
	var req productDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	p, err := rpc.ProductDetail(ctx, productID, req.Token)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"product":     productToJSON(p),
	}
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

	resp := map[string]interface{}{
		"status_code":  int64(0),
		"status_msg":   "success",
		"product_list": productsToJSON(products),
		"total":        total,
	}
	pack.SendResponse(c, resp)
}

// SeckillCreate .
// @router /seckill/activity/create/ [POST]
func SeckillCreate(ctx context.Context, c *app.RequestContext) {
	var req seckillCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}
	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activity, err := rpc.SeckillCreate(ctx, &product.CreateSeckillRequest{
		Token:        req.Token,
		ProductId:    productID,
		SeckillPrice: req.SeckillPrice,
		TotalStock:   req.TotalStock,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"activity":    seckillToJSON(activity),
	}
	pack.SendResponse(c, resp)
}

// SeckillUpdate .
// @router /seckill/activity/update/ [POST]
func SeckillUpdate(ctx context.Context, c *app.RequestContext) {
	var req seckillUpdateRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}
	activityID, err := strconv.ParseInt(req.ActivityID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}
	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activity, err := rpc.SeckillUpdate(ctx, &product.UpdateSeckillRequest{
		Token:        req.Token,
		ActivityId:   activityID,
		ProductId:    productID,
		SeckillPrice: req.SeckillPrice,
		TotalStock:   req.TotalStock,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"activity":    seckillToJSON(activity),
	}
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

	resp := map[string]interface{}{
		"status_code":   int64(0),
		"status_msg":    "success",
		"activity_list": seckillListToJSON(activities),
		"total":         total,
	}
	pack.SendResponse(c, resp)
}

// SeckillDetail .
// @router /seckill/activity/detail/ [GET]
func SeckillDetail(ctx context.Context, c *app.RequestContext) {
	var req seckillDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}
	activityID, err := strconv.ParseInt(req.ActivityID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activity, err := rpc.SeckillDetail(ctx, activityID, req.Token)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"activity":    seckillToJSON(activity),
	}
	pack.SendResponse(c, resp)
}

// ProductUpdate .
// @router /seckill/product/update/ [POST]
func ProductUpdate(ctx context.Context, c *app.RequestContext) {
	var req productUpdateRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	p, err := rpc.ProductUpdate(ctx, &product.UpdateProductRequest{
		Token:       req.Token,
		ProductId:   productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageURL,
		Category:    req.Category,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"product":     productToJSON(p),
	}

	pack.SendResponse(c, resp)
}
