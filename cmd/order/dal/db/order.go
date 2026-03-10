package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ozline/tiktok/pkg/constants"
	"github.com/ozline/tiktok/pkg/errno"
	"gorm.io/gorm"
)

type Order struct {
	Id         int64
	OrderNo    string
	UserId     int64
	ProductId  int64
	ActivityId int64
	Amount     int64
	Status     int64          `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	order.Id = SF.NextVal()
	order.OrderNo = fmt.Sprintf("SK%d", order.Id)
	err := DB.WithContext(ctx).Table(constants.OrderTableName).Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetOrderByNo(ctx context.Context, orderNo string) (*Order, error) {
	o := new(Order)
	err := DB.WithContext(ctx).Table(constants.OrderTableName).Where("order_no = ?", orderNo).First(o).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errno.OrderNotFoundError
		}
		return nil, err
	}
	return o, nil
}

func GetOrderByUserAndActivity(ctx context.Context, userId, activityId int64) (*Order, error) {
	o := new(Order)
	err := DB.WithContext(ctx).Table(constants.OrderTableName).
		Where("user_id = ? AND activity_id = ? AND status != ?", userId, activityId, constants.OrderStatusCancelled).
		First(o).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return o, nil
}

func ListOrders(ctx context.Context, userId int64, status int64, page int64, size int64) ([]*Order, int64, error) {
	var orders []*Order
	var total int64

	query := DB.WithContext(ctx).Table(constants.OrderTableName).Where("user_id = ?", userId)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func UpdateOrderStatus(ctx context.Context, orderNo string, status int64) error {
	return DB.WithContext(ctx).Table(constants.OrderTableName).
		Where("order_no = ?", orderNo).
		Update("status", status).Error
}
