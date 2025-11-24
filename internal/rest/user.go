package rest

import (
	"context"
	"myGreenMarket/domain"
	"myGreenMarket/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	Register(ctx context.Context, user *domain.User) (domain.User, error)
	Login(ctx context.Context, email, password string) (string, domain.User, error)
	VerifyEmail(ctx context.Context, verificationCodeEncrypt string) (err error)
}

type UserHandler struct {
	userService UserService
	validator   *validator.Validate
	timeout     time.Duration
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
		timeout:     10 * time.Second,
	}
}

type UserRegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var reqUser UserRegisterRequest

	if err := c.Bind(&reqUser); err != nil {
		logger.Error("Invalid request body", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	if err := h.validator.Struct(&reqUser); err != nil {
		logger.Error("Failed to validation user register", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	user, err := h.userService.Register(ctx, &domain.User{
		FullName: reqUser.FullName,
		Email:    reqUser.Email,
		Password: reqUser.Password,
	})
	if err != nil {
		logger.Error("Failed to register user", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    user,
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var reqUser UserLoginRequest

	if err := c.Bind(&reqUser); err != nil {
		logger.Error("Failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	if err := h.validator.Struct(&reqUser); err != nil {
		logger.Error("Failed to validate user login", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	token, user, err := h.userService.Login(ctx, reqUser.Email, reqUser.Password)
	if err != nil {
		logger.Error("Failed to login with user", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

func (h *UserHandler) VerifyEmail(c echo.Context) error {
	encCode := c.Param("code")

	ctx, cancel := context.WithTimeout(c.Request().Context(), h.timeout)
	defer cancel()

	err := h.userService.VerifyEmail(ctx, encCode)
	if err != nil {
		if strings.Contains(err.Error(), "invalid or expired") {
			return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, "Successfully verified email")
}
