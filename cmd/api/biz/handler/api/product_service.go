package api

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	productdb "github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/utils"
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

type productPublishRequest struct {
	Token     string `json:"token" form:"token" query:"token"`
	ProductID string `json:"product_id" form:"product_id" query:"product_id"`
}

type productDeleteRequest struct {
	Token     string `json:"token" form:"token" query:"token"`
	ProductID string `json:"product_id" form:"product_id" query:"product_id"`
}

var initProductDBOnce sync.Once

func ensureProductDB() {
	initProductDBOnce.Do(func() {
		productdb.Init()
	})
}

// ProductPublish .
// @router /seckill/product/publish/ [POST]
// publish = soft-delete-unset + set status=1
func ProductPublish(ctx context.Context, c *app.RequestContext) {
	ensureProductDB()

	var req productPublishRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// Unscoped to allow publish even if the record is soft-deleted.
	var p productdb.Product
	if err := productdb.DB.Unscoped().Where("id = ? AND seller_id = ?", productID, claims.UserId).First(&p).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// 软删除：删除时写 deleted_at；发布时清空 deleted_at。
	if err := productdb.DB.Unscoped().
		Model(&productdb.Product{}).
		Where("id = ? AND seller_id = ?", productID, claims.UserId).
		Updates(map[string]interface{}{
			"status":     1,
			"deleted_at": nil,
			"updated_at": time.Now(),
		}).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"product_id":  productID,
	}
	pack.SendResponse(c, resp)
}

// ProductDelete .
// @router /seckill/product/delete/ [POST]
// soft delete = set deleted_at
func ProductDelete(ctx context.Context, c *app.RequestContext) {
	ensureProductDB()

	var req productDeleteRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	productID, err := strconv.ParseInt(req.ProductID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// Unscoped to check deleted_at.
	var p productdb.Product
	if err := productdb.DB.Unscoped().Where("id = ? AND seller_id = ?", productID, claims.UserId).First(&p).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	if p.DeletedAt.Valid {
		// Already soft-deleted
		resp := map[string]interface{}{
			"status_code": int64(0),
			"status_msg":  "success",
			"product_id":  productID,
			"deleted_at":  p.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		}
		pack.SendResponse(c, resp)
		return
	}

	// GORM soft delete will populate deleted_at.
	if err := productdb.DB.Where("id = ? AND seller_id = ?", productID, claims.UserId).Delete(&productdb.Product{}).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// Reload to return the actual deleted_at timestamp.
	var deleted productdb.Product
	if err := productdb.DB.Unscoped().Where("id = ? AND seller_id = ?", productID, claims.UserId).First(&deleted).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"product_id":  productID,
		"deleted_at":  deleted.DeletedAt.Time.Format("2006-01-02 15:04:05"),
	}
	pack.SendResponse(c, resp)
}

type seckillDeleteRequest struct {
	Token      string `json:"token" form:"token" query:"token"`
	ActivityID string `json:"activity_id" form:"activity_id" query:"activity_id"`
}

// SeckillActivityDelete .
// @router /seckill/activity/delete/ [POST]
// soft delete = set deleted_at
func SeckillActivityDelete(ctx context.Context, c *app.RequestContext) {
	ensureProductDB()

	var req seckillDeleteRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	claims, err := utils.CheckToken(req.Token)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	activityID, err := strconv.ParseInt(req.ActivityID, 10, 64)
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// Ensure ownership: seckill_activity has no seller_id, so join product.seller_id.
	productIDs := productdb.DB.Model(&productdb.Product{}).
		Select("id").
		Where("seller_id = ?", claims.UserId)

	// Unscoped to check deleted_at.
	var act productdb.SeckillActivity
	if err := productdb.DB.Unscoped().
		Where("id = ? AND product_id IN (?)", activityID, productIDs).
		First(&act).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	if act.DeletedAt.Valid {
		resp := map[string]interface{}{
			"status_code": int64(0),
			"status_msg":  "success",
			"activity_id": activityID,
			"deleted_at":  act.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		}
		pack.SendResponse(c, resp)
		return
	}

	// GORM soft delete will populate deleted_at.
	if err := productdb.DB.Where("id = ? AND product_id IN (?)", activityID, productIDs).
		Delete(&productdb.SeckillActivity{}).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	// Reload to return the actual deleted_at timestamp.
	var deleted productdb.SeckillActivity
	if err := productdb.DB.Unscoped().
		Where("id = ? AND product_id IN (?)", activityID, productIDs).
		First(&deleted).Error; err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := map[string]interface{}{
		"status_code": int64(0),
		"status_msg":  "success",
		"activity_id": activityID,
		"deleted_at":  deleted.DeletedAt.Time.Format("2006-01-02 15:04:05"),
	}
	pack.SendResponse(c, resp)
}
