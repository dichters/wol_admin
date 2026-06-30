package middleware

const (
	TraceIdKey  = "traceId"
	SpanIdKey   = "spanId"
	HeaderTrace = "X-B3-TraceId"
	HeaderSpan  = "X-B3-SpanId"
)

// ToString 将任意类型转换为字符串
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
