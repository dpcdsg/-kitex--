package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/cmd/product/dal/cache"
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

	updated, err := db.UpdateProduct(s.ctx, p)
	if err != nil {
		return nil, err
	}

	// 让 Redis 的商品库存与 DB 对齐，否则秒杀 Lua 会基于旧库存判断。
	_ = cache.WarmUpProductStock(s.ctx, updated.Id, updated.Stock)
	return updated, nil
}
