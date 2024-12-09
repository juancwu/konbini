package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"konbini/common"
	"konbini/handler"
	"konbini/middleware"
	"konbini/store"
)

func main() {
	zerolog.TimeFieldFormat = common.FRIENDLY_TIME_FORMAT
	if os.Getenv("APP_ENV") == common.DEVELOPMENT_ENV {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: common.FRIENDLY_TIME_FORMAT})
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Failed to load .env")
		}
	}

	db, err := store.NewConn()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load .env")
	}
	defer db.Close()

	e := echo.New()

	e.HTTPErrorHandler = handler.ErrorHandler

	e.HideBanner = os.Getenv("APP_ENV") == common.PRODUCTION_ENV

	e.Use(emw.RecoverWithConfig(emw.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error().Err(err).Bytes("stack", stack).Msg("Oh uh... this is a recover attempt.")
			return err
		},
	}))
	e.Use(emw.RequestID())
	e.Use(middleware.Logger())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start TLS server.")
	}
}
