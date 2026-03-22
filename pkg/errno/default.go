package errno

var (
	Success = NewErrNo(SuccessCode, "Success")

	ServiceError             = NewErrNo(ServiceErrorCode, "service is unable to start successfully")
	ServiceInternalError     = NewErrNo(ServiceErrorCode, "service internal error")
	ParamError               = NewErrNo(ParamErrorCode, "parameter error")
	AuthorizationFailedError = NewErrNo(AuthorizationFailedErrCode, "authorization failed")
	UnexpectedTypeError      = NewErrNo(UnexpectedTypeErrorCode, "unexpected type")

	UserExistedError = NewErrNo(ParamErrorCode, "user existed")
	UserNotFoundError = NewErrNo(ParamErrorCode, "user not found")

	StockSoldOutError      = NewErrNo(StockSoldOutErrorCode, "stock sold out")
	SeckillNotStartedError = NewErrNo(SeckillNotStartedErrorCode, "seckill activity not started")
	SeckillEndedError      = NewErrNo(SeckillEndedErrorCode, "seckill activity ended")
	RepeatSeckillError     = NewErrNo(RepeatSeckillErrorCode, "repeated seckill request")

	OrderNotFoundError = NewErrNo(OrderNotFoundErrorCode, "order not found")
	OrderStatusError   = NewErrNo(OrderStatusErrorCode, "order status error")
	OrderExpiredError  = NewErrNo(OrderExpiredErrorCode, "order expired")

	ProductNotFoundError       = NewErrNo(ProductNotFoundErrorCode, "product not found")
	ProductPermissionDeniedError = NewErrNo(ProductPermissionDeniedCode, "no permission to modify this product")
)
