package service

import "context"

type ProductService struct {
	ctx context.Context
}

func NewProductService(ctx context.Context) *ProductService {
	return &ProductService{ctx: ctx}
}
