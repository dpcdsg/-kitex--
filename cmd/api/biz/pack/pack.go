package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/ozline/tiktok/pkg/errno"
)

type Response struct {
	Code int64  `json:"status_code"`
	Msg  string `json:"status_msg"`
}

func SendFailResponse(c *app.RequestContext, err error) {
	if err == nil {
		c.JSON(consts.StatusOK, Response{
			Code: errno.SuccessCode,
			Msg:  errno.SuccessMsg,
		})
		return
	}

	c.JSON(consts.StatusOK, Response{
		Code: -1,
		// 前端只展示“错误含义”，不再直接显示错误码（例如 [10103] ...）。
		Msg: errno.ConvertErr(err).ErrorMsg,
	})
}

func SendResponse(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, data)
}
