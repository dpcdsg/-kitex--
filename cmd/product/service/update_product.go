package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func (s *ProductService) UpdateProduct(req *product.UpdateProductRequest) (*db.Product, error) {
	p := &db.Product{
		Id:          req.ProductId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageUrl:    req.ImageUrl,
		Category:    req.Category,
	}

	return db.UpdateProduct(s.ctx, p)
}
