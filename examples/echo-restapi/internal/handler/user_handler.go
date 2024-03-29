package handler

import (
	"fmt"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/echokit"
	user "github.com/adipurnama/go-toolkit/examples/echo-restapi/internal"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/internal/service"
	"github.com/adipurnama/go-toolkit/examples/echo-restapi/pkg/dto"
	"github.com/adipurnama/go-toolkit/tracer"
)

// CreateUser for create single user based on json body.
func CreateUser(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(c.Request().Context(), tracer.SpanLvlServiceLogic)
		defer span.End()

		c.SetRequest(c.Request().WithContext(ctx))

		var req dto.CreateUserRequest

		err := c.Bind(&req)
		if err != nil {
			return errors.Wrap(err, "parse request body failed")
		}

		if err := echokit.Validate(c, req); err != nil {
			return err
		}

		u := user.User{
			Name:  req.Name,
			Email: req.Email,
		}

		err = svc.CreateUser(c.Request().Context(), &u)
		if err != nil {
			return err
		}

		code := http.StatusAccepted

		resp := dto.SuccessResponse{
			Status: http.StatusText(code),
			Code:   code,
			Data: dto.CreateUserResponse{
				Name:  req.Name,
				Email: req.Email,
				ID:    u.ID,
			},
		}

		return c.JSON(code, resp)
	}
}

// GetUser find user by specific id
// GET /users/:id.
func GetUser(svc *service.Service) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			msg := fmt.Sprintf("invalid user id '%s'", ctx.Param("id"))

			return echo.NewHTTPError(http.StatusBadRequest, msg).
				SetInternal(err)
		}

		user, err := svc.FindUserByID(ctx.Request().Context(), id)
		if err != nil {
			return err
		}

		resp := dto.CreateUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}

		return ctx.JSON(http.StatusOK, resp)
	}
}
