package router

import (
	"myGreenMarket/internal/rest"

	"github.com/labstack/echo/v4"
)

func SetupUserRoutes(api *echo.Group, handler *rest.UserHandler) {
	users := api.Group("/users")

	users.GET("/email-verification/:code", handler.VerifyEmail)
	users.POST("/register", handler.Register)
	users.POST("/login", handler.Login)
}
