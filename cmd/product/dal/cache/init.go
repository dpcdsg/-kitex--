package cache

import (
	"context"
	"fmt"

	"github.com/ozline/tiktok/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

const (
	SeckillStockKeyPrefix   = "seckill:stock:"
	ProductStockKeyPrefix   = "product:stock:"
	SeckillProductKeyPrefix = "seckill:product:"
)

// 同时扣减秒杀库存与商品真实库存，避免“秒杀超卖但商品库存未变化”的问题。
var DeductSeckillAndProductStockScript = redis.NewScript(`
	local seckillStock = tonumber(redis.call('get', KEYS[1]))
	local productStock = tonumber(redis.call('get', KEYS[2]))
	if seckillStock and seckillStock > 0 and productStock and productStock > 0 then
		redis.call('decr', KEYS[1])
		redis.call('decr', KEYS[2])
		return 1
	end
	return 0
`)

var RevertSeckillAndProductStockScript = redis.NewScript(`
	redis.call('incr', KEYS[1])
	redis.call('incr', KEYS[2])
	return 1
`)

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       4,
	})
}

func WarmUpProductStock(ctx context.Context, productId int64, stock int64) error {
	key := fmt.Sprintf("%s%d", ProductStockKeyPrefix, productId)
	return RDB.Set(ctx, key, stock, 0).Err()
}

func WarmUpSeckillStock(ctx context.Context, activityId int64, productId int64, stock int64) error {
	seckillKey := fmt.Sprintf("%s%d", SeckillStockKeyPrefix, activityId)
	// seckill -> product 映射，避免扣减时每次再查 DB
	seckillProductKey := fmt.Sprintf("%s%d", SeckillProductKeyPrefix, activityId)
	if err := RDB.Set(ctx, seckillKey, stock, 0).Err(); err != nil {
		return err
	}
	return RDB.Set(ctx, seckillProductKey, productId, 0).Err()
}

func GetSeckillProductID(ctx context.Context, activityId int64) (int64, error) {
	key := fmt.Sprintf("%s%d", SeckillProductKeyPrefix, activityId)
	return RDB.Get(ctx, key).Int64()
}

func DeductSeckillStock(ctx context.Context, activityId int64) (bool, int64, error) {
	seckillKey := fmt.Sprintf("%s%d", SeckillStockKeyPrefix, activityId)
	productId, err := GetSeckillProductID(ctx, activityId)
	if err != nil {
		return false, 0, err
	}
	productKey := fmt.Sprintf("%s%d", ProductStockKeyPrefix, productId)
	result, err := DeductSeckillAndProductStockScript.Run(ctx, RDB, []string{seckillKey, productKey}).Int64()
	if err != nil {
		return false, 0, err
	}
	return result == 1, productId, nil
}

func RevertSeckillAndProductStock(ctx context.Context, activityId int64, productId int64) error {
	seckillKey := fmt.Sprintf("%s%d", SeckillStockKeyPrefix, activityId)
	productKey := fmt.Sprintf("%s%d", ProductStockKeyPrefix, productId)
	_, err := RevertSeckillAndProductStockScript.Run(ctx, RDB, []string{seckillKey, productKey}).Int64()
	return err
}
