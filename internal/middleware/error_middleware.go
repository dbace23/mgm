package middleware

import (
	"myGreenMarket/pkg/logger"
	jsonres "myGreenMarket/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"

	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		if msg, ok := httpError.Message.(string); ok {
			message = msg
		}
	}

	requestID := ""
	if reqID := c.Get("requestID"); reqID != nil {
		requestID = reqID.(string)
	}

	logger.Error("Request error",
		"path", c.Request().URL.Path,
		"method", c.Request().Method,
		"error", err.Error(),
		"request_id", requestID,
	)

	if !c.Response().Committed {
		c.JSON(code, jsonres.ErrorWithRequestID("Error", message, nil, requestID))
	}
}
