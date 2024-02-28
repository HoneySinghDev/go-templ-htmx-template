package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	auth_sessions_key        = "access_token"
	auth_key          string = "authenticated"
	user_id_key       string = "user_id"
	email_key         string = "email"
)

func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(auth_sessions_key, c)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if email, ok := sess.Values[email_key].(string); ok && email != "" {
			c.Set(email_key, email) // set the email in the context
		}

		return next(c)
	}
}
