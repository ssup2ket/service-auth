package http_server

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/errors"
)

// Error response
type errResponse struct {
	HTTPStatusCode int `json:"-"`
	ErrorInfo
}

func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// Renderer
func getErrRendererBadRequest() render.Renderer {
	return &errResponse{
		ErrorInfo: ErrorInfo{
			Code:    errors.CodeBadRequest,
			Message: errors.MsgBadRequest,
		},
		HTTPStatusCode: http.StatusBadRequest, // 400
	}
}

func getErrRendererUnauthorized() render.Renderer {
	return &errResponse{
		ErrorInfo: ErrorInfo{
			Code:    errors.CodeUnauthorized,
			Message: errors.MsgUnauthorized,
		},
		HTTPStatusCode: http.StatusUnauthorized, // 401
	}
}

func getErrRendererNotFound(res errors.ErrResouce) render.Renderer {
	errCode := errors.CodeNotFound
	errMsg := errors.MsgNotFound
	switch res {
	case errors.ErrResouceUser:
		errCode = errors.CodeNotFoundUser
		errMsg = errors.MsgNotFoundUser
	}

	return &errResponse{
		ErrorInfo: ErrorInfo{
			Code:    errCode,
			Message: errMsg,
		},
		HTTPStatusCode: http.StatusNotFound, // 404
	}
}

func getErrRendererConflict(res errors.ErrResouce) render.Renderer {
	errCode := errors.CodeNotFound
	errMsg := errors.MsgNotFound
	switch res {
	case errors.ErrResouceUser:
		errCode = errors.CodeConflictUser
	}

	return &errResponse{
		ErrorInfo: ErrorInfo{
			Code:    errCode,
			Message: errMsg,
		},
		HTTPStatusCode: http.StatusConflict, // 409
	}
}

func getErrRendererServerError() render.Renderer {
	return &errResponse{
		ErrorInfo: ErrorInfo{
			Code: errors.CodeServerError,
		},
		HTTPStatusCode: http.StatusInternalServerError, // 500
	}
}
