package pack

import (
	"errors"

	"github.com/ozline/tiktok/kitex_gen/product"
	"github.com/ozline/tiktok/pkg/errno"
)

func BuildBaseResp(err error) *product.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}

	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceError.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *product.BaseResp {
	return &product.BaseResp{
		Code: err.ErrorCode,
		Msg:  err.ErrorMsg,
	}
}
