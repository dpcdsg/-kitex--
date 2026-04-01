package pack

import (
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func Product(p *product.Product) *api.Product {
	if p == nil {
		return nil
	}
	return &api.Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageURL:    p.ImageUrl,
		Category:    p.Category,
		SellerID:    0,
	}
}

func ProductList(products []*product.Product) []*api.Product {
	result := make([]*api.Product, 0, len(products))
	for _, p := range products {
		result = append(result, Product(p))
	}
	return result
}

func SeckillActivity(a *product.SeckillActivity) *api.SeckillActivity {
	if a == nil {
		return nil
	}
	return &api.SeckillActivity{
		ID:             a.Id,
		ProductID:      a.ProductId,
		ProductName:    a.ProductName,
		SeckillPrice:   a.SeckillPrice,
		TotalStock:     a.TotalStock,
		AvailableStock: a.AvailableStock,
		StartTime:      a.StartTime,
		EndTime:        a.EndTime,
		Status:         a.Status,
	}
}

func SeckillActivityList(activities []*product.SeckillActivity) []*api.SeckillActivity {
	result := make([]*api.SeckillActivity, 0, len(activities))
	for _, a := range activities {
		result = append(result, SeckillActivity(a))
	}
	return result
}
