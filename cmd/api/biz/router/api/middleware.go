package api

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ozline/tiktok/cmd/api/biz/middleware"
)

func rootMw() []app.HandlerFunc {
	return nil
}

func _seckillMw() []app.HandlerFunc {
	return nil
}

func _actionMw() []app.HandlerFunc {
	return nil
}

func _seckillactionMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _activityMw() []app.HandlerFunc {
	return nil
}

func _createMw() []app.HandlerFunc {
	return nil
}

func _seckillcreateMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _detailMw() []app.HandlerFunc {
	return nil
}

func _seckilldetailMw() []app.HandlerFunc {
	return nil
}

func _listMw() []app.HandlerFunc {
	return nil
}

func _seckilllistMw() []app.HandlerFunc {
	return nil
}

func _orderMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _cancelMw() []app.HandlerFunc {
	return nil
}

func _ordercancelMw() []app.HandlerFunc {
	return nil
}

func _detail0Mw() []app.HandlerFunc {
	return nil
}

func _orderdetailMw() []app.HandlerFunc {
	return nil
}

func _list0Mw() []app.HandlerFunc {
	return nil
}

func _orderlistMw() []app.HandlerFunc {
	return nil
}

func _payMw() []app.HandlerFunc {
	return nil
}

func _orderpayMw() []app.HandlerFunc {
	return nil
}

func _productMw() []app.HandlerFunc {
	return nil
}

func _create0Mw() []app.HandlerFunc {
	return nil
}

func _productcreateMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _detail1Mw() []app.HandlerFunc {
	return nil
}

func _productdetailMw() []app.HandlerFunc {
	return nil
}

func _list1Mw() []app.HandlerFunc {
	return nil
}

func _productlistMw() []app.HandlerFunc {
	return nil
}

func _updateMw() []app.HandlerFunc {
	return nil
}

func _productupdateMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _userMw() []app.HandlerFunc {
	return nil
}

func _userinfoMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.AuthToken()}
}

func _loginMw() []app.HandlerFunc {
	return nil
}

func _userloginMw() []app.HandlerFunc {
	return nil
}

func _userregisterMw() []app.HandlerFunc {
	return nil
}

func _registerMw() []app.HandlerFunc {
	return nil
}
