package service

import (
	"github.com/ozline/tiktok/cmd/product/dal/db"
	"github.com/ozline/tiktok/kitex_gen/product"
)

func (s *ProductService) GetSeckill(activityId int64) (*product.SeckillActivity, error) {
	activity, err := db.GetSeckillByID(s.ctx, activityId)
	if err != nil {
		return nil, err
	}

	p, err := db.GetProductByID(s.ctx, activity.ProductId)
	productName := ""
	if err == nil {
		productName = p.Name
	}

	return &product.SeckillActivity{
		Id:             activity.Id,
		ProductId:      activity.ProductId,
		ProductName:    productName,
		SeckillPrice:   activity.SeckillPrice,
		TotalStock:     activity.TotalStock,
		AvailableStock: activity.AvailableStock,
		StartTime:      activity.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:        activity.EndTime.Format("2006-01-02 15:04:05"),
		Status:         activity.Status,
	}, nil
}
