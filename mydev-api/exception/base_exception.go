package exception

// BaseException 基础异常类
type BaseException struct {
	Code int
	Msg  string
}

// Error 实现 error 接口
func (e *BaseException) Error() string {
	return e.Msg
}

// NewBaseException 创建基础异常
func NewBaseException(code int, msg string) *BaseException {
	return &BaseException{
		Code: code,
		Msg:  msg,
	}
}
