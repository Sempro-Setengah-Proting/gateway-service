package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func TestJWTAuthUserIDClaim(t *testing.T) {
	const userID = "550e8400-e29b-41d4-a716-446655440000"
	for _, tt := range []struct {
		name   string
		claims jwt.MapClaims
		key    string
		want   int
	}{
		{"identity token without sub", jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(time.Hour).Unix()}, "secret", http.StatusNoContent},
		{"missing user_id", jwt.MapClaims{"sub": userID, "exp": time.Now().Add(time.Hour).Unix()}, "secret", http.StatusUnauthorized},
		{"invalid user_id", jwt.MapClaims{"user_id": "invalid", "exp": time.Now().Add(time.Hour).Unix()}, "secret", http.StatusUnauthorized},
		{"expired", jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(-time.Hour).Unix()}, "secret", http.StatusUnauthorized},
		{"missing expiration", jwt.MapClaims{"user_id": userID}, "secret", http.StatusUnauthorized},
		{"wrong signature", jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(time.Hour).Unix()}, "wrong-secret", http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.claims).SignedString([]byte(tt.key))
			if err != nil {
				t.Fatal(err)
			}
			e := echo.New()
			e.GET("/identity/user/me", func(c *echo.Context) error {
				if c.Get("user_id") != userID {
					t.Error("incorrect user ID in context")
				}
				return c.NoContent(http.StatusNoContent)
			}, JWTAuth([]byte("secret"), nil))
			req := httptest.NewRequest(http.MethodGet, "/identity/user/me", nil)
			req.Header.Set("Authorization", "Bearer "+signed)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Fatalf("expected %d, got %d: %s", tt.want, rec.Code, rec.Body.String())
			}
		})
	}
}
