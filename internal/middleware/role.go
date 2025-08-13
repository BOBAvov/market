package middleware

import "github.com/gofiber/fiber/v2"

func RequireSeller() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals(CtxUserRole).(string)
		if role != "seller" {
			return fiber.NewError(fiber.StatusForbidden, "seller role required")
		}
		return c.Next()
	}
}
