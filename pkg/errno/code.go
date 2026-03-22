package errno

const (
	StatusSuccessCode = 0
	StatusSuccessMsg  = "ok"

	SuccessCode = 10000
	SuccessMsg  = "ok"

	ServiceErrorCode            = 10001
	ParamErrorCode              = 10002
	AuthorizationFailedErrCode  = 10003
	UnexpectedTypeErrorCode     = 10004

	StockSoldOutErrorCode       = 10101
	SeckillNotStartedErrorCode  = 10102
	SeckillEndedErrorCode       = 10103
	RepeatSeckillErrorCode      = 10104

	OrderNotFoundErrorCode      = 10201
	OrderStatusErrorCode        = 10202
	OrderExpiredErrorCode       = 10203

	ProductNotFoundErrorCode    = 10301
	ProductPermissionDeniedCode = 10302
)
