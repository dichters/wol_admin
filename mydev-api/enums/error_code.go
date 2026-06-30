package enums

// ErrorCode 错误码枚举
type ErrorCode int

const (
	// Success 成功
	Success ErrorCode = 0

	// InvalidParameter 无效参数
	InvalidParameter ErrorCode = 40001

	// DivideByZero 除零错误
	DivideByZero ErrorCode = 40002

	// UnknownError 未知错误
	UnknownError ErrorCode = 99999

	// DatabaseError 数据库错误
	DatabaseError ErrorCode = 50001
)

// String 返回错误码的字符串表示
func (e ErrorCode) String() string {
	switch e {
	case Success:
		return "Success"
	case InvalidParameter:
		return "InvalidParameter"
	case DivideByZero:
		return "DivideByZero"
	case UnknownError:
		return "UnknownError"
	case DatabaseError:
		return "DatabaseError"
	default:
		return "Unknown"
	}
}
