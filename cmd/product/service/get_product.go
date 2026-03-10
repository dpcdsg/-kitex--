package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
)

func (s *ProductService) GetProduct(productId int64) (*db.Product, error) {
	return db.GetProductByID(s.ctx, productId)
}
