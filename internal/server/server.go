package server

import (
	"context"
	"fmt"
	"log"

	"gatewayservice/internal/infrastructure/config"
	"gatewayservice/internal/router"

	"github.com/labstack/echo/v5"
	globalMiddleware "github.com/labstack/echo/v5/middleware"
)

func Run() error {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Validate configuration
	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

    if cfg.GatewayServiceKey == "" {
        return fmt.Errorf("GATEWAY_SERVICE_KEY is required")
    }

    // 3. Initialize Echo
    e := echo.New()

    // 4. Global middleware
    e.Use(globalMiddleware.Recover())

    // 5. Register routes
    if err := router.Register(e, *cfg); err != nil {
        return fmt.Errorf("register routes: %w", err)
    }
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c *echo.Context) error {
        log.Printf(
            "[IDENTITY] Incoming request: %s %s",
            c.Request().Method,
            c.Request().URL.Path,
        )

        return next(c)
    }
})

    // 6. Start server
    sc := echo.StartConfig{
        Address: ":"+cfg.GatewayPort,
    }

    return sc.Start(context.Background(), e)
}