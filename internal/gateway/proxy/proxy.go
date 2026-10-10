package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
)

type userIDKey struct{}

func New(
	targetURL string,
	prefix string,
	serviceKey string,
) (echo.HandlerFunc, error) {

	target, err := url.Parse(targetURL)
	if err != nil || target.Host == "" ||
		(target.Scheme != "http" && target.Scheme != "https") {
		return nil, fmt.Errorf("invalid upstream URL: %s", targetURL)
	}

	reverseProxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)

			// Hilangkan prefix service dari URL.
			path := strings.TrimPrefix(r.In.URL.Path, prefix)
			if path == "" {
				path = "/"
			}

			r.Out.URL.Path = path
			r.Out.URL.RawPath = ""

			// Jangan percaya identity header dari client.
			r.Out.Header.Del("X-User-ID")
			r.Out.Header.Del("X-Service-Key")

			// User ID hanya boleh berasal dari Gateway.
			if userID, ok := r.In.Context().Value(userIDKey{}).(string); ok && userID != "" {
				r.Out.Header.Set("X-User-ID", userID)
			}

			if serviceKey != "" {
				r.Out.Header.Set("X-Service-Key", serviceKey)
			}

			r.SetXForwarded()
		},

		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "upstream service unavailable", http.StatusBadGateway)
		},
	}

	return func(c *echo.Context) error {
		ctx := context.WithValue(
			c.Request().Context(),
			userIDKey{},
			c.Get("user_id"),
		)

		reverseProxy.ServeHTTP(
			c.Response(),
			c.Request().WithContext(ctx),
		)

		return nil
	}, nil
}
