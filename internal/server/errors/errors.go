package errors

// Error code
const (
	// Resource
	codeResouceUser = "_USER"

	// Common error
	CodeBadRequest   = "BAD_REQEUEST"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeServerError  = "INTERNAL_SERVER_ERROR"

	// Resource not found
	CodeNotFound     = "NOT_FOUND"
	CodeNotFoundUser = CodeNotFound + codeResouceUser

	// Resource confilct
	CodeConflict     = "CONFLICT"
	CodeConflictUser = CodeConflict + codeResouceUser
)

// Error resource
type ErrResouce string

const (
	ErrResouceUser ErrResouce = "USER"
)
