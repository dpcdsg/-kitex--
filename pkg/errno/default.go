package errno

var (
	Success = NewErrNo(SuccessCode, "Success")

	ServiceError             = NewErrNo(ServiceErrorCode, "服务启动失败")
	ServiceInternalError     = NewErrNo(ServiceErrorCode, "服务内部错误")
	ParamError               = NewErrNo(ParamErrorCode, "参数错误")
	AuthorizationFailedError = NewErrNo(AuthorizationFailedErrCode, "授权失败")
	UnexpectedTypeError      = NewErrNo(UnexpectedTypeErrorCode, "非预期类型")

	UserExistedError = NewErrNo(ParamErrorCode, "用户已存在")
	UserNotFoundError = NewErrNo(ParamErrorCode, "用户不存在")

	StockSoldOutError      = NewErrNo(StockSoldOutErrorCode, "库存不足/已售罄")
	SeckillNotStartedError = NewErrNo(SeckillNotStartedErrorCode, "秒杀尚未开始")
	SeckillEndedError      = NewErrNo(SeckillEndedErrorCode, "秒杀活动已结束")
	RepeatSeckillError     = NewErrNo(RepeatSeckillErrorCode, "重复秒杀请求")

	OrderNotFoundError = NewErrNo(OrderNotFoundErrorCode, "订单不存在")
	OrderStatusError   = NewErrNo(OrderStatusErrorCode, "订单状态错误")
	OrderExpiredError  = NewErrNo(OrderExpiredErrorCode, "订单已过期")

	ProductNotFoundError       = NewErrNo(ProductNotFoundErrorCode, "商品不存在")
	ProductPermissionDeniedError = NewErrNo(ProductPermissionDeniedCode, "无权限修改该商品")
)
