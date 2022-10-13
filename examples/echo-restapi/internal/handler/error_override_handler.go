package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	errFromService = errors.New("error while processing logic")
	errInternal    = errors.New("error while processing logic")
	statusCode400  = 400
	statusCode500  = 500
)

// EmitError helps to build echo non 2xx error response
// GET /error/:code.
func EmitError(c echo.Context) error {
	codeStr := c.Param("code")

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return err
	}

	err = CallService(c.Request().Context(), code)
	if err != nil {
		if errors.Is(err, errFromService) {
			return echo.NewHTTPError(code, fmt.Sprintf("error %d returned", code)).SetInternal(err)
		}

		return err
	}

	return c.String(http.StatusOK, "All is well")
}

// CallService - simulate service / repo layer error.
func CallService(_ context.Context, code int) error {
	if code >= statusCode400 {
		// e.g. err := query.GetDataByID(...)
		// return err
		if code > statusCode500 {
			return errors.WithStack(errInternal)
		}

		return errors.Wrapf(errFromService, "code=%d", code)
	}

	return nil
}
