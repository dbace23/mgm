package middleware

import (
	"myGreenMarket/pkg/logger"
	"myGreenMarket/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsonres "myGreenMarket/pkg/response"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, jsonres.Error(
					"UNAUTHORIZED", "Missing authorization header", nil,
				))
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, jsonres.Error(
					"UNAUTHORIZED", "Invalid authorization format", nil,
				))
			}

			claims, err := utils.ParseJWT(tokenParts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, jsonres.Error(
					"UNAUTHORIZED", "Invalid token", nil,
				))
			}

			expAt, err := claims.GetExpirationTime()
			if err != nil {
				return c.JSON(http.StatusForbidden, jsonres.Error(
					"FORBIDDEN", "Status Forbidden", nil,
				))
			}

			if time.Now().After(expAt.Time) {
				return c.JSON(http.StatusForbidden, jsonres.Error(
					"FORBIDDEN", "Status Forbidden", nil,
				))
			}

			userIDUint, err := strconv.ParseUint(claims.UserID, 10, 64)
			if err != nil {
				logger.Error("Invalid user ID in token", err)
				return c.JSON(http.StatusForbidden, jsonres.Error(
					"FORBIDDEN", "Invalid user ID in token", nil,
				))
			}

			c.Set("user_id", uint(userIDUint))
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			roleStr, ok := role.(string)
			if !ok || strings.ToUpper(roleStr) != "ADMIN" {
				return c.JSON(http.StatusForbidden, jsonres.Error(
					"FORBIDDEN", "Admin access required", nil,
				))
			}

			return next(c)
		}
	}
}
