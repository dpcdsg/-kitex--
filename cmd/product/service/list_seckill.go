package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func (s *ProductService) ListSeckill(status, page, size int64) ([]*product.SeckillActivity, int64, error) {
	activities, total, err := db.ListSeckills(s.ctx, status, page, size)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*product.SeckillActivity, 0, len(activities))
	for _, a := range activities {
		p, _ := db.GetProductByID(s.ctx, a.ProductId)
		productName := ""
		if p != nil {
			productName = p.Name
		}

		result = append(result, &product.SeckillActivity{
			Id:             a.Id,
			ProductId:      a.ProductId,
			ProductName:    productName,
			SeckillPrice:   a.SeckillPrice,
			TotalStock:     a.TotalStock,
			AvailableStock: a.AvailableStock,
			StartTime:      a.StartTime.Format("2006-01-02 15:04:05"),
			EndTime:        a.EndTime.Format("2006-01-02 15:04:05"),
			Status:         a.Status,
		})
	}

	return result, total, nil
}
