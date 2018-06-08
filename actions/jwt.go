package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/mitchellh/mapstructure"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"log"
)

func UserMiddleware() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

			fmt.Println("Hello world")
			t := c.Request().Header.Get("Authorization")
			if t != "" {
				//TODO Actualy verify the token and sign the secret properly refernced in https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html
				token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				//TODO Check the token claims to make sure that they are valid
				if claims, ok := token.Claims.(jwt.MapClaims); ok {

					log.Println("I'm here s")
					cognitoData := models.CognitoData{}
					mapstructure.Decode(claims, &cognitoData)
					tx, err := pop.Connect("development")
					if err != nil {
						log.Panic(err)
					}
					if !ok {
						return errors.WithStack(errors.New("no transaction found"))
					}

					// Allocate an empty User
					user := &models.User{}

					err = tx.Where("email = ?", cognitoData.Email).First(user)

					if err != nil {

					}

					if user.Email == "" {
						user.Name = cognitoData.CognitoUsername
						user.Email = cognitoData.Email

						err = tx.Save(user)

						if err != nil {
							return errors.WithStack(err)
						}
					}

					c.Set("user", user)
				}
			}
			return next(c)
		}
	}
}
