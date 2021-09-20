package header

// Request ID key
type ctxKeyRequestID int

const (
	RequestIDHeader = "X-Request-ID"

	RequestIDKey ctxKeyRequestID = 0
)
