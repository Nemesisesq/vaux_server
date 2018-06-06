package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/nemesisesq/vaux_server/chat"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_vaux_server_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)`
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)

		}
		app.Use(T.Middleware())
		app.Use(UserMiddleware())

		app.GET("/", HomeHandler)

		app.GET("/connect", chat.Connect)

		app.Resource("/sounds", SoundsResource{})
		app.Resource("/users", UsersResource{})
		app.Resource("/threads", ThreadsResource{})
		app.Resource("/message_sounds", MessageSoundsResource{})
		app.Resource("/messages", MessagesResource{})

		
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

func UserMiddleware() buffalo.MiddlewareFunc {
	return func (next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {

			t := c.Request().Header.Get("Authorization")
			//TODO Actualy verify the token and sign the secret properly refernced in https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-with-identity-providers.html
			token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("secret"), nil
			})
			//TODO Check the token claims to make sure that they are valid
			if claims, ok := token.Claims.(jwt.MapClaims); ok {

				fmt.Println(claims)
				//mapstructure.Decode(claims, &cognitoData)
				//Rehydrate User

				//user.Find(bson.M{"cognito_data.sub": cognitoData.Sub})
				//user.CognitoData = cognitoData
				//user.Upsert(bson.M{"user_info.sub": cognitoData.Sub})
				//return user, nil
				//json.NewEncoder(w).Encode(user)
			} else {
				//json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				//exp = Exception{Message: "Invalid authorization token"}
				//return user, exp
			}
			return next(c)
		}
	}
}