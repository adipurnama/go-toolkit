package handler

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	user "github.com/adipurnama/go-toolkit/examples/echo-restapi/internal"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/repository"
	"github.com/adipurnama/go-toolkit/log"
)

type errResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

// ErrorHandler is fallback exception handler when controller / handler returning error instead of ctx.JSON / echo.NewHTTPError.
func ErrorHandler(err error, ctx echo.Context) {
	var errEchoHTTP *echo.HTTPError

	if errors.As(err, &errEchoHTTP) {
		resp := &errResponse{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    "my custom data",
		}
		_ = ctx.JSON(resp.Code, resp)
		log.FromCtx(ctx.Request().Context()).Debug("error is echo err", "error", err)
		return
	}

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
