package v1

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Response represents a generic API response
type Response struct {
	Message string  `json:"message" example:"User created"`
	Error   *string `json:"error,omitempty" example:"failed"`
}

func NewResponse(msg string) Response {
	return Response{
		Message: msg,
	}
}

func SuccessHTTP(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, Response{
		Message: message,
	})
}

func ErrorHTTP(c echo.Context, statusCode int, err error) error {
	msg := err.Error()
	return c.JSON(statusCode, Response{
		Message: "An error occurred",
		Error:   &msg,
	})
}
