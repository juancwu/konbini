package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/juancwu/konbini/server/database"
	_ "github.com/juancwu/konbini/server/env"
	"github.com/juancwu/konbini/server/middleware"
	"github.com/juancwu/konbini/server/router"
	"github.com/juancwu/konbini/server/utils"
)

type ReqValidator struct {
	validator *validator.Validate
}

func (rq *ReqValidator) Validate(i interface{}) error {
	if err := rq.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func main() {
	/*
	   Environemnt variables are loaded in the env package when it is imported
	*/

	database.Connect()
	database.Migrate()

	e := echo.New()
	e.Use(middleware.RequestID(32))
	e.Use(middleware.Logger())
	validate := validator.New()
	validate.RegisterValidation("ValidateStringSlice", utils.ValidateStringSlice)
	e.Validator = &ReqValidator{validator: validate}

	apiV1Group := e.Group("/api/v1")
	apiV1Group.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Konbini is healthy running version: '%s' (request id: %s)", os.Getenv("APP_VERSION"), c.Request().Header.Get(echo.HeaderXRequestID)))
	})

	router.SetupAccountRoutes(apiV1Group)
	router.SetupBentoRoutes(apiV1Group)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
