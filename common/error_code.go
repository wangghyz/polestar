package common

// errorCode 错误码
type errorCode struct {
	code int
	message string
}

var (
	// HTTP相关
	ERR_HTTP_REQUEST_ERROR = &errorCode{1000, "请求错误"}
	ERR_HTTP_AUTH_FAILED = &errorCode{1001, "认证错误"}

	// 业务相关
	ERR_BUSINESS_ERROR = &errorCode{2000, "业务错误"}

	// 系统异常相关
	ERR_SYS_ERROR = &errorCode{9000, "系统错误"}
)