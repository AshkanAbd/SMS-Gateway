package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	pkgLog "github.com/AshkanAbd/arvancloud_sms_gateway/pkg/logger"
)

// HealthCheck returns application health check
//
//	@Summary		Application health check
//	@Description	Application health check
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	any
//	@Router			/healthz [get]
func HealthCheck(errChs []<-chan error) fiber.Handler {
	healthy := true
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			for i := range errChs {
				select {
				case err := <-errChs[i]:
					healthy = false
					pkgLog.Error(err, "Error in health check")
				default:
					continue
				}
			}
		}
	}()

	return func(c *fiber.Ctx) error {
		if healthy {
			return c.Status(200).Send([]byte("healthy\n"))
		} else {
			return c.Status(503).Send([]byte("unhealthy\n"))
		}
	}
}
