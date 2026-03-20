package main

import (
	"backend/dto"
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/connection"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/util"
	"log/slog"
	"net/http"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	recoverMid "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	util.SetupLogger()

	// Load config
	cnf := config.Get()

	// Connect DB
	dbConnection := connection.GetDatabase(cnf.Database)

	// Fiber App
	app := fiber.New()
	app.Use(recoverMid.New())
	app.Use(middleware.RequestLog())

	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("Unauthorized"))
		},
	})

	// ================================
	// REPOSITORIES
	// ================================
	userRepository := repository.NewUser(dbConnection)
	orderRepository := repository.NewOrder(dbConnection)
	orderItemRepository := repository.NewOrderItem(dbConnection)
	marketplaceTokenRepository := repository.NewMarketplaceToken(dbConnection)
	marketplaceRepository := repository.NewMarketplace(cnf, marketplaceTokenRepository)

	// ================================
	// SERVICES
	// ================================
	authService := service.NewAuth(cnf, userRepository)
	orderService := service.NewOrderService(orderRepository, orderItemRepository, marketplaceRepository)
	marketplaceService := service.NewMarketplace(marketplaceRepository)

	// ================================
	// API ROUTES
	// ================================
	api.NewAuth(app, authService)
	api.NewOrder(app, orderService, jwtMidd)
	app.Use("/api/marketplace", jwtMidd)
	api.NewMarketplace(app, marketplaceService)

	// ================================
	// RUN SERVER
	// ================================
	address := cnf.Server.Host + ":" + cnf.Server.Port
	slog.Info("starting server", "address", address)
	_ = app.Listen(address)
}
