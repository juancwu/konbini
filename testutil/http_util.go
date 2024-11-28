package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ResponseMatcher defines how to match response bodies
type ResponseMatcher interface {
	Matches(t *testing.T, got string) bool
}

// ExactMatcher matches the response exactly
type ExactMatcher struct {
	Expected string
}

func (m ExactMatcher) Matches(t *testing.T, got string) bool {
	return m.Expected == got
}

// JSONMatcher matches JSON responses with support for dynamic fields
type JSONMatcher struct {
	Expected map[string]interface{}
	// Fields to ignore during comparison (e.g., "token", "timestamp")
	IgnoreFields []string
}

func (m JSONMatcher) Matches(t *testing.T, got string) bool {
	var gotMap map[string]interface{}
	err := json.Unmarshal([]byte(got), &gotMap)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
		return false
	}

	// Remove ignored fields from comparison
	for _, field := range m.IgnoreFields {
		delete(gotMap, field)
		delete(m.Expected, field)
	}

	assert.Equal(t, m.Expected, gotMap)
	return true
}

// RegexMatcher matches response using regular expressions
type RegexMatcher struct {
	Pattern string
}

func (m RegexMatcher) Matches(t *testing.T, got string) bool {
	match, err := regexp.MatchString(m.Pattern, got)
	if err != nil {
		t.Errorf("Invalid regex pattern: %v", err)
		return false
	}
	return match
}

// HTTPTestCase defines a test case for HTTP handlers
type HTTPTestCase struct {
	Name           string
	Method         string
	Path           string
	Handler        echo.HandlerFunc
	Middlewares    []echo.MiddlewareFunc
	RequestBody    interface{}
	ExpectedStatus int
	ExpectedBody   ResponseMatcher
	SetupHeaders   func(*http.Request)
}

// RunHTTPTest executes an HTTP test case
func RunHTTPTest(t *testing.T, e *echo.Echo, tc HTTPTestCase) {
	// Set up the route
	switch tc.Method {
	case http.MethodGet:
		e.GET(tc.Path, tc.Handler, tc.Middlewares...)
	case http.MethodPost:
		e.POST(tc.Path, tc.Handler, tc.Middlewares...)
	case http.MethodPut:
		e.PUT(tc.Path, tc.Handler, tc.Middlewares...)
	case http.MethodDelete:
		e.DELETE(tc.Path, tc.Handler, tc.Middlewares...)
	default:
		t.Fatalf("Unsupported HTTP method: %s", tc.Method)
	}

	// Create request
	var req *http.Request
	if tc.RequestBody != nil {
		buf, err := json.Marshal(tc.RequestBody)
		assert.NoError(t, err)
		reader := bytes.NewReader(buf)
		req = httptest.NewRequest(tc.Method, tc.Path, reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(tc.Method, tc.Path, nil)
	}

	// Set custom headers if provided
	if tc.SetupHeaders != nil {
		tc.SetupHeaders(req)
	}

	// Execute request
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Assert response
	assert.Equal(t, tc.ExpectedStatus, rec.Code)
	if tc.ExpectedBody != nil {
		assert.True(t, tc.ExpectedBody.Matches(t, rec.Body.String()))
	}
}
