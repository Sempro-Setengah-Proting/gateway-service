package router

import (
	"gatewayservice/internal/infrastructure/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func TestIdentityPublicRoutes(t *testing.T) {
	var upstreamMethod, upstreamPath string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamMethod, upstreamPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()
	e := echo.New()
	if err := Register(e, config.Config{IdentityServiceURL: upstream.URL, JWTSecret: "test-secret"}); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/documentation"},
		{http.MethodGet, "/documentation/"},
		{http.MethodGet, "/openapi.json"},
		{http.MethodGet, "/health"},
		{http.MethodPost, "/refresh-token"},
		{http.MethodPost, "/oauth/sign-in-google"},
		{http.MethodPost, "/auth/sign-in"},
		{http.MethodPost, "/auth/sign-up"},
		{http.MethodPost, "/auth/register/otp/request"},
		{http.MethodPost, "/auth/register/otp/verify"},
		{http.MethodGet, "/auth/find-email/user@example.com"},
		{http.MethodPost, "/auth/forgot-password/reset"},
		{http.MethodPost, "/auth/forgot-password/otp/request"},
		{http.MethodPost, "/auth/forgot-password/otp/verify"},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, httptest.NewRequest(tt.method, "/identity"+tt.path, nil))
			if recorder.Code != http.StatusNoContent {
				t.Fatalf("expected upstream response, got %d: %s", recorder.Code, recorder.Body.String())
			}
			if upstreamMethod != tt.method || upstreamPath != tt.path {
				t.Fatalf("unexpected upstream request: %s %s", upstreamMethod, upstreamPath)
			}
		})
	}
	for _, path := range []string{"/auth/refresh", "/auth/find-email/", "/auth/find-email/user/extra", "/profile"} {
		recorder := httptest.NewRecorder()
		e.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/identity"+path, nil))
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("expected JWT guard for %s, got %d", path, recorder.Code)
		}
	}
}

func TestIdentityProfileRoutesForwardAuthenticatedUserID(t *testing.T) {
	const userID = "550e8400-e29b-41d4-a716-446655440000"
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-User-ID") != userID {
			t.Errorf("upstream received incorrect user ID: %q", r.Header.Get("X-User-ID"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer upstream.Close()
	e := echo.New()
	if err := Register(e, config.Config{IdentityServiceURL: upstream.URL, JWTSecret: "test-secret"}); err != nil {
		t.Fatal(err)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID, "exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ method, path string }{
		{http.MethodGet, "/identity/user/me"},
		{http.MethodPut, "/identity/user/edit"},
	} {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-User-ID", "client-spoofed-user")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("%s: expected upstream response, got %d: %s", tt.path, rec.Code, rec.Body.String())
		}
	}
}
