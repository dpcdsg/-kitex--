package pack

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func Product(p *db.Product) *product.Product {
	if p == nil {
		return nil
	}
	return &product.Product{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageUrl:    p.ImageUrl,
		Category:    p.Category,
		Status:      p.Status,
	}
}

func Products(ps []*db.Product) []*product.Product {
	result := make([]*product.Product, 0, len(ps))
	for _, p := range ps {
		result = append(result, Product(p))
	}
	return result
}

func SeckillActivity(a *db.SeckillActivity) *product.SeckillActivity {
	if a == nil {
		return nil
	}
	return &product.SeckillActivity{
		Id:             a.Id,
		ProductId:      a.ProductId,
		SeckillPrice:   a.SeckillPrice,
		TotalStock:     a.TotalStock,
		AvailableStock: a.AvailableStock,
		StartTime:      a.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:        a.EndTime.Format("2006-01-02 15:04:05"),
		Status:         a.Status,
	}
}
