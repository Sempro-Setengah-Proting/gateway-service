package router

import (
	"gatewayservice/internal/gateway/middleware"
	"gatewayservice/internal/gateway/proxy"
	"gatewayservice/internal/infrastructure/config"

	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo, cfg config.Config) error {

	publicRoutes := map[string]bool{
		"GET /identity/documentation":                     true,
		"GET /identity/documentation/":                    true,
		"GET /identity/openapi.json":                      true,
		"GET /identity/health":                            true,
		"POST /identity/refresh-token":                    true,
		"POST /identity/oauth/sign-in-google":             true,
		"POST /identity/auth/sign-in":                     true,
		"POST /identity/auth/sign-up":                     true,
		"POST /identity/auth/register/otp/request":        true,
		"POST /identity/auth/register/otp/verify":         true,
		"GET /identity/auth/find-email/:email":            true,
		"POST /identity/auth/forgot-password/reset":       true,
		"POST /identity/auth/forgot-password/otp/request": true,
		"POST /identity/auth/forgot-password/otp/verify":  true,
	}

	jwtGuard := middleware.JWTAuth(
		[]byte(cfg.JWTSecret),
		publicRoutes,
	)

	identityProxy, err := proxy.New(
		cfg.IdentityServiceURL,
		"/identity",
		cfg.GatewayServiceKey,
	)
	if err != nil {
		return err
	}

	e.Any("/identity/*", identityProxy, jwtGuard)

	return nil
}
