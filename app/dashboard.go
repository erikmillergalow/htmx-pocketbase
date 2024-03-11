package app

import (
	"github.com/labstack/echo/v5"

	"erikmillergalow/htmx-pocketbase/lib"
)

func GetDashboard(c echo.Context) error {
	return lib.Render(c, 200, Dashboard())
}