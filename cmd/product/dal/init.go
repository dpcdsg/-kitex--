package dal

import (
	"github.com/ozline/tiktok/cmd/product/dal/cache"
	"github.com/ozline/tiktok/cmd/product/dal/db"
)

func Init() {
	db.Init()
	cache.Init()
}
