package errors

import (
	"fmt"
)

type withCode struct {
	err   error // error 错误
	code  int   // 业务错误吗
	cause error // 原始错误
}

func WithCode(code int, format string, a ...any) error {
	return &withCode{
		err:   fmt.Errorf(format, a...),
		code:  code,
		cause: nil,
	}
}

// WrapC 封装error为withCode的错误类型
func WrapC(err error, code int, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, a...),
		code:  code,
		cause: err,
	}
}

func (w *withCode) Error() string {
	return fmt.Sprintf("%v", w)
}
