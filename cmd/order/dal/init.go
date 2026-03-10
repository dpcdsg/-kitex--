package dal

import (
	"github.com/ozline/tiktok/cmd/order/dal/cache"
	"github.com/ozline/tiktok/cmd/order/dal/db"
	"github.com/ozline/tiktok/cmd/order/dal/mq"
)

func Init() {
	db.Init()
	cache.Init()
	mq.Init()
}
