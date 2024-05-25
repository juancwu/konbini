package router

import (
	"net/http"

	"github.com/juancwu/konbini/server/middleware"
	bentomodel "github.com/juancwu/konbini/server/models/bento"
	usermodel "github.com/juancwu/konbini/server/models/user"
	"github.com/juancwu/konbini/server/service"
	"github.com/juancwu/konbini/server/utils"
	"github.com/labstack/echo/v4"
)

func SetupBentoRoutes(e *echo.Echo) {
	e.POST("/bento/personal/new", handleNewPersonalBento, middleware.JwtAuthMiddleware)
}

type NewPersonalBentoReqBody struct {
	Name       string `json:"name" validate:"required,min=1"`
	PublickKey string `json:"public_key" validate:"required,min=1"`
	Content    string `json:"content" validate:"required"`
}

func handleNewPersonalBento(c echo.Context) error {
	reqBody := new(NewPersonalBentoReqBody)

	if err := c.Bind(reqBody); err != nil {
		utils.Logger().Errorf("Failed to bind request body: %v\n", err)
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}
	if err := c.Validate(reqBody); err != nil {
		utils.Logger().Errorf("Request body validation failed: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad request")
	}

	claims := c.Get("claims").(*service.JwtCustomClaims)

	// get user
	isReal, err := usermodel.IsRealUser(claims.UserId)
	if err != nil {
		utils.Logger().Errorf("Failed to get user: %v\n", err)
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}
	if !isReal {
		utils.Logger().Error("User with id in claims returned as not real. Possible old active access token used.")
		return c.String(http.StatusBadRequest, "You must be an existing member of Konbini to create personal bentos.")
	}

	// check if user has the same personal bento
	exists, err := bentomodel.PersonalBentoExistsWithName(claims.UserId, reqBody.Name)
	if err != nil {
		utils.Logger().Errorf("Failed to check if user has personal bento with same name: %v\n", err)
		return c.String(http.StatusInternalServerError, "Personal bento service down. Please try again later.")
	}

	if exists {
		utils.Logger().Error("Attempt to create a new personal bento with the same name.")
		return c.String(http.StatusBadRequest, "Another personal bento with the same name already exists. If you wish to replace the bento, please delete it and create a new one.")
	}

	bentoId, err := bentomodel.NewPersonalBento(claims.UserId, reqBody.Name, reqBody.PublickKey, reqBody.Content)
	if err != nil {
		utils.Logger().Errorf("Failed to create new personal bento: %v\n", err)
		return c.String(http.StatusInternalServerError, "Failed to create personal bento. Please try again later.")
	}
	utils.Logger().Info("New personal bento created.", "user_id", claims.UserId, "bento_id", bentoId)

	return c.String(http.StatusCreated, bentoId)
}