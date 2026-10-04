package router

import (
	"gatewayservice/internal/gateway/handlers"

	"github.com/labstack/echo/v5"
)

type Dependencies struct {
	HealthHandler *handlers.HealthHandler
}

func Register(e *echo.Echo, deps Dependencies) {
	e.GET("/health", deps.HealthHandler.Check)
}
