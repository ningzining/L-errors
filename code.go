package errors

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	unknownCoder = defaultCoder{1, http.StatusInternalServerError, "服务器出现未知错误类型"}
)

type Coder interface {
	HTTPStatus() int

	String() string

	Code() int
}

type defaultCoder struct {
	C int

	HTTP int

	Msg string
}

func (coder defaultCoder) Code() int {
	return coder.C

}

func (coder defaultCoder) String() string {
	return coder.Msg
}

func (coder defaultCoder) HTTPStatus() int {
	if coder.HTTP == 0 {
		return http.StatusInternalServerError
	}

	return coder.HTTP
}

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

// Register 注册用户自定义的错误码
// 相同的错误码会被覆盖
func Register(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is unknownCode error code, can not be used")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	codes[coder.Code()] = coder
}

// MustRegister 注册用户自定义的错误码，建议使用
// 如果出现相同的错误码则会panic中止程序
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("code `0` is unknownCode error code, can not be used")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder 解析错误类型为Coder
// 默认错误解析为未知错误异常
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codes[v.code]; ok {
			return coder
		}
	}

	return unknownCoder
}

func init() {
	codes[unknownCoder.Code()] = unknownCoder
}
