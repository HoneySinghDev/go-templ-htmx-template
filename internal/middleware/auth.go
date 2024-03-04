package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	authSessionsKey        = "access_token"
	authKey         string = "authenticated"
	userIDKey       string = "user_id"
	emailKey        string = "email"
)

func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(authSessionsKey, c)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if auth, ok := sess.Values[authKey].(bool); !ok || !auth {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if userID, ok := sess.Values[userIDKey].(int); ok && userID != 0 {
			c.Set(userIDKey, userID) // set the user_id in the context
		}

		if email, ok := sess.Values[emailKey].(string); ok && email != "" {
			c.Set(emailKey, email) // set the email in the context
		}

		return next(c)
	}
}
