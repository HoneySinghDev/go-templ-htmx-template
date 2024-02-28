package dashboard

import (
	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	util "github.com/HoneySinghDev/go-templ-htmx-template/pkg/utils"
	"github.com/HoneySinghDev/go-templ-htmx-template/views/dashboard"
	"github.com/HoneySinghDev/go-templ-htmx-template/views/layout"
	"github.com/labstack/echo/v4"
)

func HandleDashboard(_ *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		return util.Render(c, dashboard.HomeIndex(layout.PageInfo{}, dashboard.Home()))
	}
}
