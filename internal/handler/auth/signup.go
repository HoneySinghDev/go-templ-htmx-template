package auth

import (
	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	util "github.com/HoneySinghDev/go-templ-htmx-template/pkg/utils"
	viewauth "github.com/HoneySinghDev/go-templ-htmx-template/views/auth"
	"github.com/HoneySinghDev/go-templ-htmx-template/views/layout"
)

func HandleSignupIndex(_ *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return util.Render(c, viewauth.SignupIndex(layout.PageInfo{
			Title:       "Signup",
			Description: "Signup for a new account",
		}, viewauth.SignupPage()))
	}
}

func HandleSignupCreate(_ *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !htmx.IsHTMX(c.Request()) {
			return c.String(http.StatusBadRequest, "")
		}

		email := c.FormValue("email")
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirmPassword")

		u := UserCreds{
			EmailID:         email,
			Password:        password,
			ConfirmPassword: confirmPassword,
		}

		if errors, ok := u.Validate(); !ok {
			var errs string
			for _, v := range errors {
				errs += v[0] + "\n"
			}
			return util.Render(c, viewauth.SignupForm(viewauth.SignupData{
				Email: email,
				Error: errs,
			}))
		}

		// Here you would typically check if the email is already used
		// and if not, create the user in your database.
		// Skipping detailed implementation as it depends on your user management system.

		return util.Render(c, viewauth.SignupSuccess(u.EmailID))
	}
}
