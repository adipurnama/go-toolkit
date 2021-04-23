package controller

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// Todo ...
type Todo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

// CreateTodo ...
func CreateTodo(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, `{"message": "Created"}`)
}

// UpdateTodo ...
func UpdateTodo(ctx echo.Context) error {
	return ctx.JSON(http.StatusAccepted, `{"message": "Updated"}`)
}

// DeleteTodo ...
func DeleteTodo(ctx echo.Context) error {
	return ctx.JSON(http.StatusAccepted, `{"message": "Updated"}`)
}

// ListTodo ...
func ListTodo(ctx echo.Context) error {
	type response struct {
		Data []Todo `json:"data"`
	}

	var resp response

	return ctx.JSON(http.StatusAccepted, resp)
}
