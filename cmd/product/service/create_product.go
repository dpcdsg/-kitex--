package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func (s *ProductService) CreateProduct(req *product.CreateProductRequest) (*db.Product, error) {
	p := &db.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageUrl,
		Category:    req.Category,
		Status:      1,
	}

	return db.CreateProduct(s.ctx, p)
}
