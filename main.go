package main

import (
	"gatewayservice/internal/gateway/handlers"
	"gatewayservice/internal/infrastructure/config"
	"gatewayservice/internal/router"
	"log"

	"github.com/labstack/echo/v5"
)

func main() {
	cfg := config.LoadConfig()
	e := echo.New()

	healthHandler := handlers.NewHealthHandler()
	router.Register(e, router.Dependencies{
		HealthHandler: healthHandler,
	})

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
