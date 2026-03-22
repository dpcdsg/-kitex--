package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/errno"
)

func (s *ProductService) UpdateProduct(userId int64, req *product.UpdateProductRequest) (*db.Product, error) {
	if req.ProductId <= 0 || len(req.Name) == 0 || req.Price <= 0 || req.Stock < 0 {
		return nil, errno.ParamError
	}

	p, err := db.GetProductByID(s.ctx, req.ProductId)
	if err != nil {
		return nil, errno.ProductNotFoundError
	}
	if p.SellerId != userId {
		return nil, errno.ProductPermissionDeniedError
	}

	p.Name = req.Name
	p.Description = req.Description
	p.Price = req.Price
	p.Stock = req.Stock
	p.ImageUrl = req.ImageUrl
	p.Category = req.Category

	if err := db.UpdateProduct(s.ctx, p); err != nil {
		return nil, err
	}
	return db.GetProductByID(s.ctx, req.ProductId)
}
