package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type accessTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTAuth(
	secret []byte,
	publicRoutes map[string]bool,
) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()

			routeKey := req.Method + " " + req.URL.Path

			// Endpoint publik tidak perlu access token.
			if isPublicRoute(routeKey, publicRoutes) {
				return next(c)
			}

			authorization := req.Header.Get("Authorization")
			parts := strings.Fields(authorization)

			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "missing or invalid bearer token",
				})
			}

			claims := &accessTokenClaims{}

			token, err := jwt.ParseWithClaims(
				parts[1],
				claims,
				func(t *jwt.Token) (any, error) {
					return secret, nil
				},
				jwt.WithValidMethods([]string{"HS256"}),
				jwt.WithExpirationRequired(),
			)

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid or expired token",
				})
			}

			if _, err := uuid.Parse(claims.UserID); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "missing or invalid token user_id",
				})
			}

			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}

// Match named path parameters as exactly one non-empty segment.
func isPublicRoute(routeKey string, publicRoutes map[string]bool) bool {
	if publicRoutes[routeKey] {
		return true
	}
	actual := strings.Split(routeKey, "/")
	for pattern, public := range publicRoutes {
		if !public {
			continue
		}
		expected := strings.Split(pattern, "/")
		if len(actual) != len(expected) {
			continue
		}
		matched := true
		for i, segment := range expected {
			if strings.HasPrefix(segment, ":") && actual[i] != "" {
				continue
			}
			if segment != actual[i] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}
