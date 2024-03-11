package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"erikmillergalow/htmx-pocketbase/lib"
	"erikmillergalow/htmx-pocketbase/middleware"
)

type LoginFormValue struct {
	username string
	password string
}

type SignUpFormValue struct {
	username string
	password string
	passwordConfirm string
	email string
}

func getLoginFormValue(c echo.Context) LoginFormValue {
	return LoginFormValue {
		username: c.FormValue("username"),
		password: c.FormValue("password"),
	}
}

func getSignUpFormValue(c echo.Context) SignUpFormValue {
	return SignUpFormValue {
		username: c.FormValue("username"),
		password: c.FormValue("password"),
		passwordConfirm: c.FormValue("passwordConfirm"),
		email: c.FormValue("email"),
	}
}

func RegisterLoginRoutes(e *core.ServeEvent, group echo.Group) {

	// sign in page
	group.GET("/sign-in", func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/dashboard")
		}

		return lib.Render(c, 200, LoginForm(LoginFormValue{}, nil))
	})

	group.POST("/sign-in", func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/dashboard")
		}

		form := getLoginFormValue(c)

		fmt.Println(form)
		// add validation here

		var token *string
		token, err := lib.Login(e, form.username, form.password)

		if err != nil {
			fmt.Println("Error logging in")
			return lib.Render(c, 200, LoginForm(form, err))
		}

		c.SetCookie(&http.Cookie{
			Name: middleware.AuthCookieName,
			Value: *token,
			Path: "/",
			Secure: true,
			HttpOnly: true,
		})

		return lib.HtmxRedirect(c, "/app/dashboard")
	})

	group.GET("/sign-up", func(c echo.Context) error {
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/dashboard")
		}

		return lib.Render(c, 200, SignUpForm(SignUpFormValue{}, nil))
	})

	group.POST("/sign-up", func(c echo.Context) error {
		// https://github.com/pocketbase/pocketbase/discussions/4409
		if c.Get(apis.ContextAuthRecordKey) != nil {
			return c.Redirect(302, "/dashboard")
		}

		form := getSignUpFormValue(c)

		fmt.Println(form)

		var token *string
		token, err := lib.SignUp(e, c, form.username, form.password, form.passwordConfirm)

		if err != nil {
			fmt.Println("Error signing up")
			return lib.Render(c, 200, SignUpForm(form, err))
		}

		c.SetCookie(&http.Cookie{
			Name: middleware.AuthCookieName,
			Value: *token,
			Path: "/",
			Secure: true,
			HttpOnly: true,
		})

		return lib.HtmxRedirect(c, "/app/dashboard")
	})

	group.POST("/logout", func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name: middleware.AuthCookieName,
			Value: "",
			Path: "/",
			Secure: true,
			HttpOnly: true,
			MaxAge: -1,
		})

		return lib.HtmxRedirect(c, "/auth/sign-in")
	})
}