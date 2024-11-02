package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
)

type AuthHandler struct {
	service   ports.AuthService
	validator *validator.Validate
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		service:   authService,
		validator: validator.New(),
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	credentials := new(domain.UserCredentials)

	if err := c.BodyParser(credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request body
	if err := h.validator.Struct(credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}
	// Attempt login
	user, token, err := h.service.Login(credentials)
	if err != nil {
		// Handle different types of errors appropriately
		switch err.Error() {
		case "invalid credentials":
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "An error occurred during login",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token required",
		})
	}

	// Remove "Bearer " prefix if present
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	tokens, err := h.service.RefreshAccessToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	return c.JSON(tokens)
}
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	request := new(domain.User)

	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request body
	if err := h.validator.Struct(request); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = formatValidationError(err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	// Attempt registration
	token, err := h.service.Register(request)
	if err != nil {
		switch err.Error() {
		case "email already registered":
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		default:
			fmt.Println("An error occurred during registration:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "An error occurred during registration",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		// "user":  user,
		"token": token,
	})
}
