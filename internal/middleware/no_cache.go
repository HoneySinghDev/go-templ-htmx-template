package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

type NoCacheConfig struct {
	Skipper middleware.Skipper
}

// NewDefaultNoCacheConfig returns a new instance of NoCacheConfig with default settings.
func NewDefaultNoCacheConfig() NoCacheConfig {
	return NoCacheConfig{
		Skipper: middleware.DefaultSkipper,
	}
}

// NoCache returns a middleware that sets several HTTP headers to prevent responses from being cached.
func NoCache() echo.MiddlewareFunc {
	return NoCacheWithConfig(NewDefaultNoCacheConfig())
}

// NoCacheWithConfig returns a NoCache middleware with custom configuration.
func NoCacheWithConfig(config NoCacheConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			epoch := time.Unix(0, 0).Format(time.RFC1123)

			noCacheHeaders := map[string]string{
				"Expires":         epoch,
				"Cache-Control":   "no-cache, private, max-age=0",
				"Pragma":          "no-cache",
				"X-Accel-Expires": "0",
			}

			etagHeaders := []string{
				"ETag",
				"If-Modified-Since",
				"If-Match",
				"If-None-Match",
				"If-Range",
				"If-Unmodified-Since",
			}

			req := c.Request()

			for _, header := range etagHeaders {
				if req.Header.Get(header) != "" {
					req.Header.Del(header)
				}
			}

			res := c.Response()
			for k, v := range noCacheHeaders {
				res.Header().Set(k, v)
			}

			return next(c)
		}
	}
}
