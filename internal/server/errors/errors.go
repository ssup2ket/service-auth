package errors

// Error code
const (
	// Code
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

	// Message
	// Resource
	msgResourcesUser = "User "

	// Common error
	MsgBadRequest   = "Bad Request"
	MsgUnauthorized = "Unauthroized"
	MsgServerError  = "Internal server error"

	// Resource not found
	MsgNotFound     = "Not found"
	MsgNotFoundUser = msgResourcesUser + MsgNotFound

	// Resource conflict
	MsgConflict     = "Conflit"
	MsgConflictUser = msgResourcesUser + MsgConflict
)

// Error resource
type ErrResouce string

const (
	ErrResouceUser ErrResouce = "USER"
)
