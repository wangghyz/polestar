package util

type (
	ResponseStatus int

	ResponseMessage struct {
		// code
		Code    ResponseStatus
		// 错误消息
		Message string
		// 返回数据
		Data    interface{}
	}
)

const (
	// 正常
	ResponseStatusSuccess     ResponseStatus = 1
	// 逻辑错误
	ResponseStatusLogicError  ResponseStatus = 2
	// 系统错误
	ResponseStatusSystemError ResponseStatus = 3
)
