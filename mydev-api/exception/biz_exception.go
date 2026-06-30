package exception

// BizException 业务异常类
type BizException struct {
	*BaseException
}

// NewBizException 创建业务异常
func NewBizException(code int, msg string) *BizException {
	return &BizException{
		BaseException: NewBaseException(code, msg),
	}
}
