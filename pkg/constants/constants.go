package constants

import "time"

const (
	JWTValue = "MTAxNTkwMTg1Mw=="
	StartID  = 10000

	// Redis DB
	RedisDBProduct = 4
	RedisDBOrder   = 5

	// Redis 缓存时间
	SeckillStockExpire = 24 * time.Hour
	OrderExpire        = 30 * time.Minute
	LockTime           = 1 * time.Second
	LockWaitTime       = 5 * time.Millisecond
	MaxRetryTimes      = 3

	// Redis Key
	SeckillStockKeyPrefix = "seckill:stock:"
	SeckillOrderKeyPrefix = "seckill:order:"
	OrderKeyPrefix        = "order:"

	// RPC
	MuxConnection  = 1
	RPCTimeout     = 3 * time.Second
	ConnectTimeout = 50 * time.Millisecond

	// 服务名
	APIServiceName     = "api"
	UserServiceName    = "user"
	ProductServiceName = "product"
	OrderServiceName   = "order"

	// 数据库表名
	UserTableName     = "user"
	ProductTableName  = "product"
	SeckillTableName  = "seckill_activity"
	OrderTableName    = "order"

	// Snowflake
	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	// 限制
	MaxConnections  = 1000
	MaxQPS          = 100
	MaxListLength   = 100
	MaxIdleConns    = 10
	MaxGoroutines   = 10
	MaxOpenConns    = 100
	ConnMaxLifetime = 10 * time.Second

	// 分页默认值
	PageNum  = 1
	PageSize = 10

	// 订单状态
	OrderStatusPending   = 0
	OrderStatusPaid      = 1
	OrderStatusCancelled = 2
	OrderStatusExpired   = 3

	// 秒杀活动状态
	SeckillStatusPending = 0
	SeckillStatusActive  = 1
	SeckillStatusEnded   = 2

	// 商品状态
	ProductStatusOn  = 1
	ProductStatusOff = 0

	// MQ
	SeckillOrderQueue = "seckill_order"
)
