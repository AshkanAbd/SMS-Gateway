package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Metrics returns application metrics
//
//	@Summary		Application metrics
//	@Description	Application metrics
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	any
//	@Router			/metrics [get]
func Metrics() fiber.Handler {
	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	return func(c *fiber.Ctx) error {
		p(c.Context())
		return nil
	}
}
