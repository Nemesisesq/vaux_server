package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/pkg/errors"
	"github.com/gobuffalo/pop"
	"golang.org/x/crypto/bcrypt"
	"github.com/mitchellh/mapstructure"
)

// AuthLogin default implementation.
func AuthLogin(c buffalo.Context) error {

	userForm := &models.UserForm{}

	if err := c.Bind(userForm); err != nil {
		return errors.WithStack(err)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	err := tx.Where("username =  ?", userForm.Email).First(userForm)

	if err != nil {
		return c.Error(404, err)
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	if string(ph) != userForm.PasswordHash {
		c.Render(200, r.JSON(errors.New("password is incorrect")))
	}
	var claims map[string]interface{}
	mapstructure.Decode(userForm.User, claims)

	jwt, err := generateJwtToken(claims)

	return c.Render(200, r.JSON(map[string]interface{}{
		"status": "successfully logged in",
		"jwt":    jwt,
	}))
}

// AuthSignup default implementation.
func AuthSignup(c buffalo.Context) error {
	userForm := &models.UserForm{}

	if err := c.Bind(userForm); err != nil {
		return errors.WithStack(err)
	}

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	verrs, err := userForm.ValidateAndCreate(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.JSON(verrs))
	}
	return c.Render(200, r.JSON(map[string]interface{}{"message": "signup successful"}))
}

// AuthResetPassword default implementation.
func AuthResetPassword(c buffalo.Context) error {
	return c.Render(200, r.HTML("auth/reset_password.html"))
}
