package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"

	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		routeErr := c.Next()
		if routeErr != nil {
			if err := c.App().ErrorHandler(c, routeErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		latency := time.Since(start)
		status := c.Response().StatusCode()

		switch {
		case status < 400:
			pkgLog.Debug("%d %s %s %.3f ms", status, c.Method(), c.Path(), float64(latency.Microseconds())/1000)
		case status >= 400 && status < 500:
			pkgLog.Info("%d %s %s %.3f ms", status, c.Method(), c.Path(), float64(latency.Microseconds())/1000)
		case status >= 500:
			pkgLog.Warn("%d %s %s %.3f ms", status, c.Method(), c.Path(), float64(latency.Microseconds())/1000)
		}

		return nil
	}
}
