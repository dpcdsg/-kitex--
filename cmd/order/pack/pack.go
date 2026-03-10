package pack

import (
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/kitex_gen/order"
	"github.com/ozline/tiktok/pkg/errno"
)

func BuildBaseResp(err error) *order.BaseResp {
	if err == nil {
		return &order.BaseResp{Code: errno.SuccessCode, Msg: errno.SuccessMsg}
	}
	e := errno.ConvertErr(err)
	return &order.BaseResp{Code: e.ErrorCode, Msg: e.ErrorMsg}
}

func Order(data *db.Order) *order.Order {
	if data == nil {
		return nil
	}
	return &order.Order{
		Id:         data.Id,
		OrderNo:    data.OrderNo,
		UserId:     data.UserId,
		ProductId:  data.ProductId,
		ActivityId: data.ActivityId,
		Amount:     data.Amount,
		Status:     data.Status,
		CreatedAt:  data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func Orders(data []*db.Order) []*order.Order {
	result := make([]*order.Order, 0, len(data))
	for _, d := range data {
		result = append(result, Order(d))
	}
	return result
}
