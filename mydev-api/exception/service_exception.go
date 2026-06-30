package exception

// ServiceException 服务异常类
type ServiceException struct {
	*BaseException
}

// NewServiceException 创建服务异常
func NewServiceException(code int, msg string) *ServiceException {
	return &ServiceException{
		BaseException: NewBaseException(code, msg),
	}
}
