package handler

import (
	"net/http"

	user "github.com/adipurnama/go-toolkit/example/echo-restapi/internal"
	"github.com/adipurnama/go-toolkit/example/echo-restapi/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type errResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ErrorHandler is fallback exception handler when controller / handler returning error instead of ctx.JSON / echo.NewHTTPError.
func ErrorHandler(err error, ctx echo.Context) {
	resp := &errResponse{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
	}

	var eNotFound user.ErrUserIDNotFound
	if ok := errors.As(err, &eNotFound); ok {
		resp.Code = http.StatusNotFound
		resp.Message = eNotFound.Error()

		_ = ctx.JSON(resp.Code, resp)

		return
	}

	if errors.Is(err, repository.ErrNoRows) {
		resp.Message = err.Error()
		resp.Code = http.StatusNotFound

		_ = ctx.JSON(resp.Code, resp)

		return
	}
}
