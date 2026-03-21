package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RequestLog() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		startedAt := time.Now()

		err := ctx.Next()

		attrs := []any{
			"method", ctx.Method(),
			"path", ctx.Path(),
			"status", ctx.Response().StatusCode(),
			"latency_ms", time.Since(startedAt).Milliseconds(),
			"ip", ctx.IP(),
			"user_agent", ctx.Get(fiber.HeaderUserAgent),
		}

		if err != nil {
			attrs = append(attrs, "error", err.Error())
			slog.Error("http request failed", attrs...)
			return err
		}

		if ctx.Response().StatusCode() >= fiber.StatusInternalServerError {
			slog.Error("http request completed with server error", attrs...)
			return nil
		}

		slog.Info("http request completed", attrs...)
		return nil
	}
}
