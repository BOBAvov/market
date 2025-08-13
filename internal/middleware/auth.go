package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"market/internal/utils"
)

type AuthConfig struct {
	JWTSecret string
}

const CtxUserID = "uid"
const CtxUserRole = "role"

func AuthRequired(cfg AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseJWT(token, cfg.JWTSecret)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals(CtxUserID, claims.UserID)
		c.Locals(CtxUserRole, claims.Role)
		return c.Next()
	}
}
