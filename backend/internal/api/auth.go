package api

import (
	"backend/domain"
	"backend/dto"
	"backend/internal/util"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type authApi struct {
	authService domain.AuthService
}

func NewAuth(app *fiber.App, authService domain.AuthService) {
	a := authApi{
		authService: authService,
	}

	app.Post("/api/auth/login", a.Login)
	app.Post("/api/auth/register", a.Register)
}

func (a *authApi) Login(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	res, err := a.authService.Login(c, req)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (a *authApi) Register(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(util.HTTPStatusFromError(util.ErrInvalidRequestBody)).JSON(dto.CreateResponseError(util.ErrInvalidRequestBody.Error()))
	}

	res, err := a.authService.Register(c, req)
	if err != nil {
		return ctx.Status(util.HTTPStatusFromError(err)).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}
