package main

import (
	"log"
	"os"

	// "erikmillergalow/htmx-pocketbase/components"

	// "github.com/a-h/templ"
	// "github.com/labstack/echo/v5"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"erikmillergalow/htmx-pocketbase/app"
	"erikmillergalow/htmx-pocketbase/auth"
	"erikmillergalow/htmx-pocketbase/middleware"
)

func main() {
	pb := pocketbase.New()

	// set the token into a cookie after successful auth
	pb.OnRecordAuthRequest().Add(func(e *core.RecordAuthEvent) error {
		e.HttpContext.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    e.Token,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
			MaxAge:   int(pb.Settings().RecordAuthToken.Duration),
			Path:     "/",
		})
		return nil
	})

	// serves static files from the provided public dir (if exists)
	pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.Static("/pb_public", "pb_public")
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		authGroup := e.Router.Group("/auth", middleware.LoadAuthContextFromCookie(pb))
		auth.RegisterLoginRoutes(e, *authGroup)

		appGroup := e.Router.Group("/app", middleware.LoadAuthContextFromCookie(pb), middleware.AuthGuard)
		appGroup.GET("/dashboard", app.GetDashboard)

		return nil
	})

	if err := pb.Start(); err != nil {
		log.Fatal(err)
	}
}
