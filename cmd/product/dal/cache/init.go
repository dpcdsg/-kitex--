package cache

import (
	"context"
	"fmt"

	"github.com/ozline/tiktok/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

const SeckillStockKeyPrefix = "seckill:stock:"

var DeductStockScript = redis.NewScript(`
	local stock = tonumber(redis.call('get', KEYS[1]))
	if stock and stock > 0 then
		redis.call('decr', KEYS[1])
		return 1
	end
	return 0
`)

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       4,
	})
}

func WarmUpSeckillStock(ctx context.Context, activityId int64, stock int64) error {
	key := fmt.Sprintf("%s%d", SeckillStockKeyPrefix, activityId)
	return RDB.Set(ctx, key, stock, 0).Err()
}

func DeductSeckillStock(ctx context.Context, activityId int64) (bool, error) {
	key := fmt.Sprintf("%s%d", SeckillStockKeyPrefix, activityId)
	result, err := DeductStockScript.Run(ctx, RDB, []string{key}).Int64()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
