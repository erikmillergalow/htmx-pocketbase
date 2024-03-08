package main

import (
    "log"
    "os"

	// "erikmillergalow/htmx-pocketbase/components"

	// "github.com/a-h/templ"
	// "github.com/labstack/echo/v5"
	// "net/http"
    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/apis"
    "github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    // serves static files from the provided public dir (if exists)
    app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

        e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
    
		// e.Router.GET("/hello/:name", func(c echo.Context) error {
		// 	name := c.PathParam("name")
		// 	component := components.Index(name)
	
		// 	return c.Render(http.StatusOK, "", component);
		// })
	
        return nil
	})

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}