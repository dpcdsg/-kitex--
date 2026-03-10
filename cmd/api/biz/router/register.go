package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/ozline/tiktok/cmd/api/biz/router/api"
)

func GeneratedRegister(r *server.Hertz) {
	api.Register(r)
}
