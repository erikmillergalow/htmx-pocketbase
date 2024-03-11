package middleware

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

const AuthCookieName = "Auth"

func LoadAuthContextFromCookie(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check for cookie
			tokenCookie, err := c.Request().Cookie(AuthCookieName)
			if err != nil {
				// no cookie
				return next(c)
			}

			token := tokenCookie.Value
			
			claims, err := security.ParseUnverifiedJWT(token)
			if err != nil {
				return next(c)
			}
			// determine token type (admin vs. user)
			tokenType := cast.ToString(claims["type"])

			switch tokenType {
			case tokens.TypeAdmin:
					admin, err := app.Dao().FindAdminByToken(
						token,
						app.Settings().AdminAuthToken.Secret,
					)
					if err == nil && admin != nil {
						fmt.Println("Setting admin by cookie")
						// set the cookie to authenticate the admin user
						c.Set(apis.ContextAdminKey, admin)
					}
			case tokens.TypeAuthRecord:
					record, err := app.Dao().FindAuthRecordByToken(
						token,
						app.Settings().RecordAuthToken.Secret,
					)
					if err == nil && record != nil {
						fmt.Println("Setting user by cookie")
						// set cookie to authenticate normal user
						c.Set(apis.ContextAuthRecordKey, record)
					}
			}
			
			return next(c)

		}
	}
}

func AuthGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		record := c.Get(apis.ContextAuthRecordKey)

		if record == nil {
			return c.Redirect(302, "auth/sign-in")
		}

		return next(c)
	}
}