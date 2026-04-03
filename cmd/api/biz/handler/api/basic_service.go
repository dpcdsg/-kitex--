package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/cmd/api/biz/middleware/metrics"
	"github.com/ozline/tiktok/cmd/api/biz/pack"
	"github.com/ozline/tiktok/cmd/api/biz/rpc"
	"github.com/ozline/tiktok/kitex_gen/user"
)

// UserRegister .
// @router /seckill/user/register/ [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {
	var req api.UserRegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.UserRegisterResponse)
	var err error
	resp.UserID, resp.Token, err = rpc.UserRegister(ctx, &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	pack.SendResponse(c, resp)
}

// UserLogin .
// @router /seckill/user/login/ [POST]
func UserLogin(ctx context.Context, c *app.RequestContext) {
	var bizErr error
	defer func() { metrics.RecordBusinessResult(metrics.HandlerUserLogin, bizErr) }()

	var req api.UserLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.UserLoginResponse)
	var err error
	resp.UserID, resp.Token, err = rpc.UserLogin(ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		bizErr = err
		pack.SendFailResponse(c, err)
		return
	}

	pack.SendResponse(c, resp)
}

// UserInfo .
// @router /seckill/user/ [GET]
func UserInfo(ctx context.Context, c *app.RequestContext) {
	var req api.UserRequest
	if err := c.BindAndValidate(&req); err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp := new(api.UserResponse)
	u, err := rpc.UserInfo(ctx, &user.InfoRequest{
		UserId: req.UserID,
		Token:  req.Token,
	})
	if err != nil {
		pack.SendFailResponse(c, err)
		return
	}

	resp.User = pack.User(u)
	pack.SendResponse(c, resp)
}
