package handlers

import (
	"gatewayservice/internal/httpresponse"
	"net/http"

	"github.com/labstack/echo/v5"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *echo.Context) error {
	return c.JSON(http.StatusOK, httpresponse.Success(
		map[string]string{"service": "gateway-service"},
		"service is healthy",
	))
}
