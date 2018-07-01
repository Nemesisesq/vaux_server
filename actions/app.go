package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/nemesisesq/vaux_server/chat"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/unrolled/secure"
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

		api := app.Group("/api")
		//api.Use(ValidateTokensFromHeader)




		app.GET("/", HomeHandler)


		app.GET("/verify", ValidateTokensFromHeader(func(c buffalo.Context) error {

			u :=  c.Value("user")
			return c.Render(http.StatusOK, r.JSON(u))
		}))

		//api endpoints auth locked
		api.GET("/connect", chat.Connect)
		api.Resource("/sounds", SoundsResource{})
		api.Resource("/users", UsersResource{})
		api.Resource("/threads", ThreadsResource{})
		api.Resource("/message_sounds", MessageSoundsResource{})
		api.Resource("/messages", MessagesResource{})

		api.Resource("/thread_members", ThreadMembersResource{})
		api.Resource("/profiles", ProfilesResource{})

		app.POST("/auth/login", AuthLogin)
		app.POST("/auth/signup", AuthSignup)
		app.POST("/auth/reset_password", AuthResetPassword)

		//public endpoints
		app.Resource("/sounds", SoundsResource{})
		app.Resource("/users", UsersResource{})
		app.Resource("/threads", ThreadsResource{})
		app.Resource("/message_sounds", MessageSoundsResource{})
		app.Resource("/messages", MessagesResource{})

		app.Resource("/thread_members", ThreadMembersResource{})
		app.Resource("/profiles", ProfilesResource{})

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}
