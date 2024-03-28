package lib

import (
	"fmt"

	"github.com/labstack/echo/v5"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

func Login(e *core.ServeEvent, username string, password string) (*string, error) {
	user, err := e.App.Dao().FindAuthRecordByUsername("users", username)
	if err != nil {
		fmt.Println("Login failed")
		fmt.Println(username)
		return nil, fmt.Errorf("Login failed")
	}

	valid := user.ValidatePassword(password)
	if !valid {
		fmt.Println("Password validation failed")

		return nil, fmt.Errorf("Password validation failed")
	}

	s, tokenErr := tokens.NewRecordAuthToken(e.App, user)
	if tokenErr != nil {
		fmt.Println("Failed to generate token")

		return nil, fmt.Errorf("Failed to generate token")
	}

	return &s, nil
}

func SignUp(e *core.ServeEvent, c echo.Context, username string, password string, passwordConfirm string) (*string, error) {
	usersCollection, err := e.App.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		fmt.Println("Failed to find user collection")
		return nil, fmt.Errorf("Failed to find user collection")
	}
	newUser := models.NewRecord(usersCollection)
	form := forms.NewRecordUpsert(e.App, newUser)

	// must use `enctype="multipart/form-data"` on form
	form.LoadRequest(c.Request(), "")
	fmt.Println(form)
	if err := form.Submit(); err != nil {
		fmt.Println("Submit form failed")
		fmt.Println(err)
		return nil, fmt.Errorf("Submit form failed")
	}

	user, err := e.App.Dao().FindAuthRecordByUsername("users", username)
	if err != nil {
		fmt.Println("Sign up failed")
		fmt.Println(username)
		return nil, fmt.Errorf("Sign up failed")
	}

	s, tokenErr := tokens.NewRecordAuthToken(e.App, user)
	if tokenErr != nil {
		fmt.Println("Failed to generate token")

		return nil, fmt.Errorf("Failed to generate token")
	}

	return &s, nil
}
