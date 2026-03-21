package middleware

import (
	"backend/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthorizeRoles(allowedRoles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.Locals("user").(*jwt.Token)
		if !ok || token == nil {
			return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("Unauthorized"))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("Unauthorized"))
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return ctx.Status(http.StatusForbidden).JSON(dto.CreateResponseError("Forbidden"))
		}

		if _, exists := allowed[role]; !exists {
			return ctx.Status(http.StatusForbidden).JSON(dto.CreateResponseError("Forbidden"))
		}

		return ctx.Next()
	}
}
