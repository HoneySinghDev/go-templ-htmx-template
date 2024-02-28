package handler

import (
	"github.com/HoneySinghDev/go-templ-htmx-template/internal/handler/auth"
	"github.com/HoneySinghDev/go-templ-htmx-template/internal/handler/dashboard"
	"github.com/HoneySinghDev/go-templ-htmx-template/internal/middleware"
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(s *server.Server) {
	routes := []*echo.Route{
		// Auth - Login
		s.Echo.GET("/login", auth.HandleLoginIndex(s)),
		s.Echo.POST("/login", auth.HandleLoginCreate(s)),
		// Auth - Signup
		s.Echo.GET("/signup", auth.HandleSignupIndex(s)),
		s.Echo.POST("/signup", auth.HandleSignupCreate(s)),
	}

	s.Router.Routes = append(s.Router.Routes, routes...)
}

func DashBoardRoute(s *server.Server) {
	d := s.Echo.Group("/dashboard", middleware.WithAuth)

	routes := []*echo.Route{
		d.GET("", dashboard.HandleDashboard(s)),
	}

	s.Router.Routes = append(s.Router.Routes, routes...)
}
