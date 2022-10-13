package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	zero  = 0
	three = 3
)

// PanicGuaranteed demonstrates what could go wrong during executing handler.
func PanicGuaranteed(c echo.Context) error {
	result := three / zero

	return c.String(http.StatusOK, fmt.Sprintf("%d", result))
}
