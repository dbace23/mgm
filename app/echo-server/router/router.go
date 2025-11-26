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

func SetupProductRoutes(api *echo.Group, handler *rest.ProductHandler) {
	products := api.Group("/products")

	products.GET("", handler.GetAllProducts)
	products.POST("", handler.CreateProduct)
	products.PUT("/:id", handler.UpdateProduct)
	products.DELETE("/:id", handler.DeleteProduct)
}
