package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	accessSecret string
}

func NewAuthMiddleware(accessSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		accessSecret: accessSecret,
	}
}

// AuthenticateJWT middleware verifies the JWT token in the Authorization header
func (m *AuthMiddleware) AuthenticateJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Check if the header starts with "Bearer "
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		tokenString := headerParts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.accessSecret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Check token validity
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// Check token expiration
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token expired",
			})
		}

		// Check token type
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token type",
			})
		}

		// Store user information in context
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// RequireRole middleware checks if the user has the required role
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user role not found in context",
			})
		}

		hasRole := false
		for _, role := range roles {
			if role == userRole.(string) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}

		return c.Next()
	}
}
