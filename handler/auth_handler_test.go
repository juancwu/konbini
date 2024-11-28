package handler

import (
	"bytes"
	"encoding/json"
	"konbini/middleware"
	"konbini/store"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandleSignUp(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
		requestBody    signUpBody
		requestMethod  string
	}{
		{
			name:           "Successful sign up",
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":201,"message":"Successfully signed up."}`,
			requestBody: signUpBody{
				Email:    "my@mail.com",
				Password: "123456789abc@",
			},
			requestMethod: http.MethodPost,
		},
	}

	e := echo.New()
	db, err := store.NewConn()
	assert.Nil(t, err)
	h := authHandler{db}
	e.POST("/", h.handleSignUp, middleware.BindAndValidate(reflect.TypeOf(signUpBody{})))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf, err := json.Marshal(test.requestBody)
			assert.Nil(t, err)

			reader := bytes.NewReader(buf)

			req := httptest.NewRequest(test.requestMethod, "/", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatus, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get(echo.HeaderContentType))
			assert.Equal(t, test.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
