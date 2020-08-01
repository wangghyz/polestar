package common

// ErrorCode 错误码
type ErrorCode struct {
	Code    int
	Message string
}

var (
	// HTTP相关
	ERR_HTTP_REQUEST_ERROR = &ErrorCode{1000, "请求错误"}
	ERR_HTTP_AUTH_FAILED = &ErrorCode{1001, "认证错误"}

	// 业务相关
	ERR_BUSINESS_ERROR = &ErrorCode{2000, "业务错误"}

	// 系统异常相关
	ERR_SYS_ERROR = &ErrorCode{9000, "系统错误"}
)