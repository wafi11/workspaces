package pkg

import (
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, status int, message string, err error) error {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return c.JSON(status, Response{
		Status:  status,
		Message: message,
		Error:   errMsg,
	})
}
