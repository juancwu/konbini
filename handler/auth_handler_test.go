package handler

import (
	"konbini/middleware"
	"konbini/store"
	"konbini/testutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers(t *testing.T) {
	e := echo.New()
	db, err := store.NewConn()
	assert.Nil(t, err)
	h := authHandler{db}

	tests := []testutil.HTTPTestCase{
		{
			Name:   "Successful sign up",
			Method: http.MethodPost,
			Path:   "/sign-up",
			RequestBody: signUpBody{
				Email:    "my@mail.com",
				Password: "123456789abc@",
			},
			Handler: h.handleSignUp,
			Middlewares: []echo.MiddlewareFunc{
				middleware.BindAndValidate(reflect.TypeOf(signUpBody{})),
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedBody: testutil.JSONMatcher{
				Expected: map[string]interface{}{
					"status":  http.StatusCreated,
					"message": "Successfully signed up.",
				},
			},
		},
		{
			Name:   "Successful sign in",
			Method: http.MethodPost,
			Path:   "/sign-in",
			RequestBody: signInBody{
				Email:    "my@mail.com",
				Password: "123456789abc@",
			},
			Handler: h.handleSignIn,
			Middlewares: []echo.MiddlewareFunc{
				middleware.BindAndValidate(reflect.TypeOf(signInBody{})),
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedBody: testutil.JSONMatcher{
				Expected: map[string]interface{}{
					"status": http.StatusOK,
				},
				IgnoreFields: []string{"token"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			testutil.RunHTTPTest(t, e, tc)
		})
	}
}
