package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/pkg/errors"
	"github.com/gobuffalo/pop"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type D struct {
	i interface{} `json:"i"`
}

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

	u := &models.User{}
	err := tx.Where("email =  ?", userForm.Email).First(u)

	if err != nil {
		return c.Error(404, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(userForm.Password))
	if err != nil {
		return errors.WithStack(err)
	}

	twoWeeks := time.Now().Add(2 * 7 * 24 * time.Hour)
	claims := getJWSClaims(*u, twoWeeks)

	refreshToken, err := generateJwtToken(claims)


	u.RefreshToken.String = string(refreshToken)

	tx.Save(u)

	twoHours := time.Now().Add(2 * time.Hour)
	claims = getJWSClaims(*u, twoHours)

	jwt, err := generateJwtToken(claims)

	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.Auto(c, map[string]interface{}{
		"status": "successfully logged in",
		"jwt":    string(jwt),
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
