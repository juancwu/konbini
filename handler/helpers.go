package handler

import (
	"github.com/labstack/echo/v4"
)

// Sends back a basic json response with a status code and message.
func respond(code int, message string, c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(code, basicResponse{Status: code, Message: message})
}
