package auth

import (
	"context"
	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"net/http"
	"time"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	util "github.com/HoneySinghDev/go-templ-htmx-template/pkg/utils"
	viewauth "github.com/HoneySinghDev/go-templ-htmx-template/views/auth"
	"github.com/HoneySinghDev/go-templ-htmx-template/views/layout"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/rs/zerolog/log"
)

const (
	auth_sessions_key        = "access_token"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	email_key         string = "email"
)

func HandleLoginIndex(_ *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return util.Render(c, viewauth.LoginIndex(layout.PageInfo{
			Title:       "Login",
			Description: "Login to your account",
		}, viewauth.LoginPage()))
	}
}

func HandleLoginCreate(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !htmx.IsHTMX(c.Request()) {
			return c.String(http.StatusBadRequest, "")
		}

		emailID := c.FormValue("email")
		password := c.FormValue("password")

		u := UserCreds{
			EmailID:  emailID,
			Password: password,
		}

		if _, ok := u.Validate(); !ok {
			return util.Render(c, viewauth.LoginForm(viewauth.LoginData{
				Email:    emailID,
				Password: "",
				Error:    "Credentials You Have Entered Are Invalid",
			}))
		}

		authUser, err := s.SB.Auth.SignIn(context.Background(), supabase.UserCredentials{
			Email:    u.EmailID,
			Password: u.Password,
		})
		if err != nil {
			log.Debug().
				Err(err).
				Str("emailID", u.EmailID).
				Msg("Error while signing in")

			return util.Render(c, viewauth.LoginForm(viewauth.LoginData{
				Email:    emailID,
				Password: "",
				Error:    "Credentails You Have Entered Are Invalid",
			}))
		}

		sess, err := session.Get(auth_sessions_key, c)
		if err != nil {
			log.Debug().
				Err(err).
				Msg("Error while getting session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   time.Now().Add(s.Config.Auth.AccessTokenValidity.GoDuration()).Second(),
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:    true,
			user_id_key: authUser.User.ID,
			email_key:   authUser.User.Email,
		}

		_ = sess.Save(c.Request(), c.Response())

		return htmx.NewResponse().Redirect("/dashboard").Write(c.Response())
	}
}
