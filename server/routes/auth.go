package routes

import (
	commonApi "konbini/common/api"
	"konbini/server/handlers"
	"konbini/server/middlewares"
	"reflect"
)

func setupAuthRoutes(routeConfig *RouteConfig) {
	routeConfig.Echo.POST(
		"/auth/login",
		handlers.Login(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.LoginRequest{})),
	)
	routeConfig.Echo.POST(
		"/auth/register",
		handlers.Register(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RegisterRequest{})),
	)

	routeConfig.Echo.POST(
		"/auth/token/check",
		handlers.CheckAuthToken(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.CheckAuthTokenRequest{})),
	)

	routeConfig.Echo.POST(
		"/auth/totp/setup",
		handlers.SetupTOTP(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
	)
	routeConfig.Echo.POST(
		"/auth/totp/lock",
		handlers.SetupTOTPLock(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(commonApi.SetupTOTPLockRequest{})),
	)
	routeConfig.Echo.DELETE(
		"/auth/totp",
		handlers.RemoveTOTP(routeConfig.DBConnector),
		middlewares.ProtectFull(routeConfig.DBConnector),
		middlewares.ValidateJson(reflect.TypeOf(handlers.RemoveTOTPRequest{})),
	)

	routeConfig.Echo.GET("/auth/email/verify", handlers.VerifyEmail(routeConfig.DBConnector))
	routeConfig.Echo.POST(
		"/auth/email/resend-verification",
		handlers.ResendVerificationEmail(routeConfig.DBConnector),
		middlewares.ProtectAll(routeConfig.DBConnector),
	)
}
