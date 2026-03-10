package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
)

func (s *ProductService) ListProducts(page, size int64, category string) ([]*db.Product, int64, error) {
	return db.ListProducts(s.ctx, page, size, category)
}
