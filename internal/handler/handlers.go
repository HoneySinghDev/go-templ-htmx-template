package handler

import (
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
)

func AttachAllRoutes(s *server.Server) {
	// attach our routes
	AuthRoutes(s)
	DashBoardRoute(s)
}
