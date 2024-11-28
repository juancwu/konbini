package handler

import (
	"database/sql"
	"konbini/middleware"
	"konbini/store"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type authHandler struct {
	db *sql.DB
}

// Setup auth routes to the group route.
// Prefix all routes with "/auth"
func SetupAuthRoutes(e *echo.Group, db *sql.DB) {
	h := authHandler{db}
	e.POST("/auth/signup", h.handleSignUp, middleware.BindAndValidate(reflect.TypeOf(signUpBody{})))
}

// Handle incoming requests for sign up.
func (h *authHandler) handleSignUp(c echo.Context) error {
	body, err := middleware.GetRequestBody[signUpBody](c)
	if err != nil {
		return err
	}

	exists, err := store.ExistsUserWithEmail(c.Request().Context(), h.db, body.Email)
	if err != nil {
		return err
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request. Email already occupied.")
	}

	userID, err := store.NewUser(c.Request().Context(), h.db, body.Email, body.Password, body.Nickname)
	if err != nil {
		return err
	}
	log.Info().Msgf("New user created with ID: %s", userID)

	return respond(http.StatusCreated, "Successfully signed up.", c)
}

func (h *authHandler) handleSignIn(c echo.Context) error {
	return respond(http.StatusOK, "", c)
}
