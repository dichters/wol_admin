# mydev-go

A Fiber-based web service template project.

## Project Structure

```
mydev-go/
├── mydev-service/    # Web backend service
├── mydev-api/        # Common types (exceptions, responses, enums)
├── mydev-sdk/        # SDK for external callers
└── log/              # Log files (gitignored)
```

## Requirements

- Go 1.25.6
- Fiber v2.50.0

## Setup

```bash
# Download dependencies
go work sync

# Run service
cd mydev-service/cmd/server
go run main.go
```

## API Endpoints

All endpoints use POST method and return JSON:

### Health Check
```
POST /srv/v1/hc
```

Response:
```json
{
  "code": 0,
  "msg": "Success.",
  "data": []
}
```

### Divide
```
POST /srv/v1/divide
Content-Type: application/json

{
  "a": 10,
  "b": 2
}
```

Response:
```json
{
  "code": 0,
  "msg": "Success.",
  "data": [5]
}
```

## SDK Usage

```go
import "github.com/mydev/mydev-sdk/client"

client := client.NewClient("http://localhost:8080")

// Health check
result, err := client.HealthCheck(traceId, spanId)

// Divide
result, err := client.Divide(10, 2, traceId, spanId)
```

## Tracing

The service uses B3 Propagation protocol for distributed tracing:
- `X-B3-TraceId`: Trace ID
- `X-B3-SpanId`: Span ID

If not provided in request headers, new IDs will be auto-generated.

## Exception Handling

- **BizException**: Business errors (returns custom code and message)
- **ServiceException**: Service errors (returns custom code, but message is masked)
- **Unknown errors**: Returns code 99999 with "Internal Server Error."

## Logging

- Default level: INFO
- Outputs to: stdout, `log/service.log`, `log/error.log`
- ERROR level and above also writes to `log/error.log`
