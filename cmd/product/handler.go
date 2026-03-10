package main

import (
	"context"

	"github.com/ozline/tiktok/cmd/product/pack"
	"github.com/ozline/tiktok/cmd/product/service"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/errno"
)

type ProductServiceImpl struct{}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (resp *product.CreateProductResponse, err error) {
	resp = new(product.CreateProductResponse)

	if len(req.Name) == 0 || req.Price <= 0 || req.Stock < 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	p, err := service.NewProductService(ctx).CreateProduct(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Product = pack.Product(p)
	return
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *product.GetProductRequest) (resp *product.GetProductResponse, err error) {
	resp = new(product.GetProductResponse)

	if req.ProductId <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	p, err := service.NewProductService(ctx).GetProduct(req.ProductId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Product = pack.Product(p)
	return
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsRequest) (resp *product.ListProductsResponse, err error) {
	resp = new(product.ListProductsResponse)

	if req.Page <= 0 || req.Size <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	products, total, err := service.NewProductService(ctx).ListProducts(req.Page, req.Size, req.Category)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.ProductList = pack.Products(products)
	resp.Total = total
	return
}

func (s *ProductServiceImpl) CreateSeckill(ctx context.Context, req *product.CreateSeckillRequest) (resp *product.CreateSeckillResponse, err error) {
	resp = new(product.CreateSeckillResponse)

	if req.ProductId <= 0 || req.SeckillPrice <= 0 || req.TotalStock <= 0 || len(req.StartTime) == 0 || len(req.EndTime) == 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	activity, err := service.NewProductService(ctx).CreateSeckill(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Activity = pack.SeckillActivity(activity)
	return
}

func (s *ProductServiceImpl) GetSeckill(ctx context.Context, req *product.GetSeckillRequest) (resp *product.GetSeckillResponse, err error) {
	resp = new(product.GetSeckillResponse)

	if req.ActivityId <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	activity, err := service.NewProductService(ctx).GetSeckill(req.ActivityId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Activity = activity
	return
}

func (s *ProductServiceImpl) ListSeckill(ctx context.Context, req *product.ListSeckillRequest) (resp *product.ListSeckillResponse, err error) {
	resp = new(product.ListSeckillResponse)

	if req.Page <= 0 || req.Size <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	activities, total, err := service.NewProductService(ctx).ListSeckill(req.Status, req.Page, req.Size)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.ActivityList = activities
	resp.Total = total
	return
}

func (s *ProductServiceImpl) DeductStock(ctx context.Context, req *product.DeductStockRequest) (resp *product.DeductStockResponse, err error) {
	resp = new(product.DeductStockResponse)

	if req.ActivityId <= 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamError)
		return resp, nil
	}

	err = service.NewProductService(ctx).DeductStock(req.ActivityId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(nil)
	return
}
