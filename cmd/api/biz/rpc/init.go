package rpc

import (
	"github.com/ozline/tiktok/kitex_gen/order/orderservice"
	"github.com/ozline/tiktok/kitex_gen/product/productservice"
	"github.com/ozline/tiktok/kitex_gen/user/userservice"
)

var (
	userClient    userservice.Client
	productClient productservice.Client
	orderClient   orderservice.Client
)

func Init() {
	InitUserRPC()
	InitProductRPC()
	InitOrderRPC()
}
