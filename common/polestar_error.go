package common

import (
	"fmt"
	"reflect"
)

// PolestarError PolestarError
type PolestarError struct {
	code          *ErrorCode
	detailMessage string
}

// Error PolestarError error实现
func (pe *PolestarError) Error() string {
	if pe.detailMessage == "" {
		return fmt.Sprintf("[%d][%s]", pe.code.Code, pe.code.Message)
	} else {
		return fmt.Sprintf("[%d][%s] %s", pe.code.Code, pe.code.Message, pe.detailMessage)
	}
}

// NewPolestarError 创建NewPolestarError
func NewPolestarError(code *ErrorCode, detailMessage string) *PolestarError {
	return &PolestarError{code, detailMessage}
}

// IsPolestarError 判断是否是PolestarError
func IsPolestarError(src interface{}) (*PolestarError, bool) {
	tp := reflect.TypeOf(src)

	if src == nil {
		return nil, false
	}
	if tp.Kind() == reflect.Ptr {
		pe, ok := src.(*PolestarError)
		return pe, ok
	} else {
		pe, ok := src.(PolestarError)
		return &pe, ok
	}
}

// PanicPolestarError 抛出PolestarError异常
func PanicPolestarError(code *ErrorCode, detailMessage string) {
	panic(NewPolestarError(code, detailMessage))
}

// PanicPolestarErrorByError 抛出PolestarError异常
func PanicPolestarErrorByError(code *ErrorCode, err error) {
	panic(NewPolestarError(code, err.Error()))
}
