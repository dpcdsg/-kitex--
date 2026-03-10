package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	api "github.com/ozline/tiktok/cmd/api/biz/handler/api"
	"github.com/ozline/tiktok/cmd/api/biz/middleware"
)

func Register(r *server.Hertz) {
	root := r.Group("/seckill")
	{
		// 用户模块
		userGroup := root.Group("/user")
		{
			userGroup.POST("/register/", api.UserRegister)
			userGroup.POST("/login/", api.UserLogin)
			userGroup.GET("/", middleware.AuthToken(), api.UserInfo)
		}

		// 商品模块
		productGroup := root.Group("/product")
		{
			productGroup.POST("/create/", middleware.AuthToken(), api.ProductCreate)
			productGroup.GET("/detail/", api.ProductDetail)
			productGroup.GET("/list/", api.ProductList)
		}

		// 秒杀活动模块
		activityGroup := root.Group("/activity")
		{
			activityGroup.POST("/create/", middleware.AuthToken(), api.SeckillCreate)
			activityGroup.GET("/list/", api.SeckillList)
			activityGroup.GET("/detail/", api.SeckillDetail)
		}

		// 秒杀下单
		root.POST("/action/", middleware.AuthToken(), api.SeckillAction)

		// 订单模块
		orderGroup := root.Group("/order", middleware.AuthToken())
		{
			orderGroup.GET("/list/", api.OrderList)
			orderGroup.GET("/detail/", api.OrderDetail)
			orderGroup.POST("/pay/", api.OrderPay)
			orderGroup.POST("/cancel/", api.OrderCancel)
		}
	}
}
