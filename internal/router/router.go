package router

import (
	"net/http"
	"strings"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/handler"

	"github.com/HoneySinghDev/go-templ-htmx-template/internal/middleware"
	"github.com/labstack/echo-contrib/session"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

//nolint:funlen
func Init(s *server.Server) {
	s.Echo = echo.New()

	s.Echo.Debug = s.Config.Echo.Debug
	s.Echo.HideBanner = true
	s.Echo.Logger.SetOutput(&echoLogger{level: s.Config.LogLevelFromString(s.Config.Logger.RequestLevel),
		log: log.With().Str("component", "echo").Logger()})

	s.Echo.HTTPErrorHandler = CustomHTTPErrorHandler

	if s.Config.Echo.RecoverMiddleware {
		s.Echo.Use(echoMiddleware.Recover())
	} else {
		log.Warn().Msg("Disabling recover middleware due to environment config")
	}

	if s.Config.Echo.SecureMiddleware.Enable {
		s.Echo.Use(echoMiddleware.SecureWithConfig(echoMiddleware.SecureConfig{
			Skipper:               echoMiddleware.DefaultSecureConfig.Skipper,
			XSSProtection:         s.Config.Echo.SecureMiddleware.XssProtection,
			ContentTypeNosniff:    s.Config.Echo.SecureMiddleware.ContentTypeNosniff,
			XFrameOptions:         s.Config.Echo.SecureMiddleware.XFrameOptions,
			HSTSMaxAge:            s.Config.Echo.SecureMiddleware.HstsMaxAge,
			HSTSExcludeSubdomains: s.Config.Echo.SecureMiddleware.HstsExcludeSubdomains,
			ContentSecurityPolicy: s.Config.Echo.SecureMiddleware.ContentSecurityPolicy,
			CSPReportOnly:         s.Config.Echo.SecureMiddleware.CspReportOnly,
			HSTSPreloadEnabled:    s.Config.Echo.SecureMiddleware.HstsPreload,
			ReferrerPolicy:        s.Config.Echo.SecureMiddleware.ReferrerPolicy,
		}))
	} else {
		log.Warn().Msg("Disabling secure middleware due to environment config")
	}

	if s.Config.Echo.LoggerMiddleware {
		s.Echo.Use(middleware.WithConfig(middleware.LoggerConfig{
			Level:             s.Config.LogLevelFromString(s.Config.Logger.RequestLevel),
			LogRequestHeader:  s.Config.Logger.RequestHeader,
			LogRequestBody:    s.Config.Logger.RequestBody,
			LogRequestQuery:   s.Config.Logger.RequestQuery,
			LogResponseHeader: s.Config.Logger.ResponseHeader,
			LogResponseBody:   s.Config.Logger.ResponseBody,
			LogCaller:         s.Config.Logger.LogCaller,
			RequestBodyLogSkipper: func(req *http.Request) bool {
				// We skip all body logging for auth endpoints as these might contain credentials
				if strings.HasPrefix(req.URL.Path, "/api/v1/auth") {
					return true
				}

				return middleware.DefaultRequestBodyLogSkipper(req)
			},
			ResponseBodyLogSkipper: func(req *http.Request, res *echo.Response) bool {
				// We skip all body logging for auth endpoints as these might contain credentials
				if strings.HasPrefix(req.URL.Path, "/api/v1/auth") {
					return true
				}

				return middleware.DefaultResponseBodyLogSkipper(req, res)
			},
			Skipper: func(c echo.Context) bool {
				// We skip logging of readiness and liveness endpoints
				switch c.Path() {
				case "/-/ready", "/-/healthy":
					return true
				}
				return false
			},
		}))
	} else {
		log.Warn().Msg("Disabling logger middleware due to environment config")
	}

	// Add your custom / additional middlewares here.
	// see https://echo.labstack.com/middleware

	// ---
	// Initialize our general groups and set middleware to use above them
	s.Router = &server.Router{
		Routes: nil, // will be populated by handlers.AttachAllRoutes(s)

		// Unsecured base group available at /**
		Root: s.Echo.Group(""),

		// Management endpoints, uncacheable, secured by key auth (query param), available at /-/**
		Management: s.Echo.Group("/-", echoMiddleware.KeyAuthWithConfig(echoMiddleware.KeyAuthConfig{
			KeyLookup: "query:mgmt-secret",
			Validator: func(key string, _ echo.Context) (bool, error) {
				return key == s.Config.Management.Secret, nil
			},
			Skipper: func(c echo.Context) bool {
				// We skip key auth for readiness and liveness endpoints
				return c.Path() == "/-/ready"
			},
		}), middleware.NoCache()),
	}

	// Session Middleware
	s.Echo.Use(session.Middleware(sessions.NewCookieStore([]byte(s.Config.Management.Secret))))

	// Static file server
	s.Echo.Static("/static", "static")

	// ---
	// Finally, attach our handlers
	handler.AttachAllRoutes(s)
}
