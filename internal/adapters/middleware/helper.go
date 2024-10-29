package middleware

import "github.com/gofiber/fiber/v2"

// ExtractUserID helper function to get user ID from context
func ExtractUserID(c *fiber.Ctx) string {
	userID := c.Locals("user_id")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// ExtractUserRole helper function to get user role from context
func ExtractUserRole(c *fiber.Ctx) string {
	role := c.Locals("role")
	if role == nil {
		return ""
	}
	return role.(string)
}
