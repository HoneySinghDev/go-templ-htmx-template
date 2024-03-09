package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"

	"github.com/HoneySinghDev/go-templ-htmx-template/pkg/server"
	util "github.com/HoneySinghDev/go-templ-htmx-template/pkg/utils"
	viewauth "github.com/HoneySinghDev/go-templ-htmx-template/views/auth"
	"github.com/HoneySinghDev/go-templ-htmx-template/views/layout"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/rs/zerolog/log"
)

const (
	authSessionsKey        = "access_token"
	authKey         string = "authenticated"
	userIDKey       string = "user_id"
	emailKey        string = "email"
	Time24Hours            = 24
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

		SaveSession(c, authUser)

		return htmx.NewResponse().Redirect("/dashboard").Write(c.Response())
	}
}

func SaveSession(c echo.Context, user *supabase.AuthenticatedDetails) {
	sess, err := session.Get(authSessionsKey, c)
	if err != nil || sess == nil {
		log.Debug().
			Err(err).
			Msg("Error while getting session")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   time.Now().Add(Time24Hours * time.Hour).Second(),
		HttpOnly: true,
	}

	// Set user as authenticated, their username,
	// their ID and the client's time zone
	sess.Values = map[interface{}]interface{}{
		authKey:   true,
		userIDKey: user.User.ID,
		emailKey:  user.User.Email,
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Debug().
			Err(err).
			Str("emailID", user.User.Email).
			Msg("Error while saving session")
	}
}
