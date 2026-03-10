package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/ozline/tiktok/config"
	"github.com/ozline/tiktok/pkg/constants"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func Init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       constants.RedisDBOrder,
	})
}

func CheckSeckillRepeat(ctx context.Context, userId, activityId int64) (bool, error) {
	key := fmt.Sprintf("%s%d:%d", constants.SeckillOrderKeyPrefix, activityId, userId)
	result, err := RDB.SetNX(ctx, key, 1, 24*time.Hour).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func RemoveSeckillRepeatFlag(ctx context.Context, userId, activityId int64) {
	key := fmt.Sprintf("%s%d:%d", constants.SeckillOrderKeyPrefix, activityId, userId)
	RDB.Del(ctx, key)
}
