package router

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/juancwu/konbini/middleware"
	"github.com/juancwu/konbini/tag"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// apiError represents a way to define a custom apiError that ErrHandler can identify and use.
type apiError struct {
	// Code is the http response status code.
	Code int `json:"-"`
	// Err is the original error. This gets log by ErrHandler.
	// This field can also be a validator.ValidationErrs.
	Err error `json:"-"`
	// Msg is the internal message that gets log.
	Msg string `json:"-"`
	// The path the api error occurred.
	Path string `json:"-"`

	// Errs is a string that gets send back to the client after ValidationErrs is processed. This field will be in the json response if not empty.
	Errs []string `json:"errors,omitempty"`
	// PublicMsg is the message that gets send back to the client. This field will be in the json response if not empty.
	PublicMsg string `json:"message,omitempty"`
	// RequestId is the id of the request the error came from, helps with traceback. This field will be in the json response if not empty.
	RequestId string `json:"request_id,omitempty"`
}

// Error satisfies the error interface so that apiError can be used as an error type.
func (e apiError) Error() string {
	return e.Msg
}

// ErrHandler is a custom error handler that will log the error and corresponding message.
// Use an echo.HTTPError if there is a need to return a status code other than 500.
// Normal errors will be handled using a 500 and generic internal server error message..
func ErrHandler(err error, c echo.Context) {
	switch err.(type) {
	case *echo.HTTPError:
		he := err.(*echo.HTTPError)
		logger := log.Info()
		if he.Code >= 400 && he.Code < 500 {
			logger = log.Warn()
		} else if he.Code >= 500 {
			logger = log.Error()
		}
		logger.Err(he).Str(echo.HeaderXRequestID, c.Request().Header.Get(echo.HeaderXRequestID)).Send()
		writeJSON(he.Code, c, basicRespBody{Msg: fmt.Sprintf("%v", he.Message), RequestId: c.Request().Header.Get(echo.HeaderXRequestID)})
	case apiError:
		e := err.(apiError)
		if ve, ok := e.Err.(validator.ValidationErrors); ok {
			if e.PublicMsg == "" {
				e.PublicMsg = "Invalid request body. Please fix the issues."
			}
			if e.Code == 0 {
				e.Code = http.StatusBadRequest
			}
			e.Errs = make([]string, len(ve))
			structType, ok := c.Get(middleware.STRUCT_TYPE_KEY).(reflect.Type)
			if !ok {
				for i, err := range ve {
					e.Errs[i] = fmt.Sprintf("Validation failed on the '%s' failed.", err.Tag())
				}
			} else {
				for i, err := range ve {
					msg := tag.ParseErrorMsgTag(structType, err)
					if msg == "" {
						msg = fmt.Sprintf("Validation failed on the '%s' failed.", err.Tag())
					}
					e.Errs[i] = msg
				}
			}
		} else if errors.Is(e.Err, middleware.ErrNoJwtClaims) {
			e.Code = http.StatusUnauthorized
			if e.Msg == "" {
				e.Msg = "Failed to get jwt claims from context."
			}
		}

		if e.PublicMsg == "" {
			e.PublicMsg = http.StatusText(e.Code)
		}

		if e.Code == 0 {
			e.Code = http.StatusInternalServerError
		}

		log.Error().
			Err(e.Err).
			Str(echo.HeaderXRequestID, e.RequestId).
			Int("status_code", e.Code).
			Str("path", e.Path).
			Msg(e.Msg)

		writeJSON(e.Code, c, e)
	default:
		log.Error().Err(err).Msg("Standard error encountered. Somewhere in route code is returning standard error.")
		c.NoContent(http.StatusInternalServerError)
	}
}
