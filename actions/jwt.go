package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func UserMiddleware() buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

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


					claims["username"] = claims["cognito:username"]

					cognitoData := models.CognitoData{}
					tmp, err := json.Marshal(claims)

					if err != nil {
						return errors.WithStack(err)
					}
					err = json.Unmarshal(tmp, &cognitoData)

					if err != nil {
						return errors.WithStack(err)
					}

					tx, ok := c.Value("tx").(*pop.Connection)
					if !ok {
						return errors.WithStack(errors.New("no transaction found"))
					}

					// Allocate an empty User
					user := &models.User{}

					err = tx.Where("email = ?", cognitoData.Email).First(user)

					if err != nil {
						user.Name = cognitoData.CognitoUsername
						user.Email = cognitoData.Email

						err = tx.Save(user)

						if err != nil {
							return errors.WithStack(err)
						}
					}
					log.Info("setting user")
					c.Set("user", *user)
				} else {
					return errors.WithStack(errors.New("The claims for this token we not valid"))
				}
			}

			return next(c)
		}
	}
}
