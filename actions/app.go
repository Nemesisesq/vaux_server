package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"
	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
	"github.com/nemesisesq/vaux_server/models"
	"github.com/nemesisesq/vaux_server/chat"
	"net/http"
	"github.com/gobuffalo/buffalo/middleware/csrf"
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

		app.GET("/verify", func (c buffalo.Context) error{
			return c.Render(http.StatusOK, r.JSON(map[string]string{"data": "Looks like everything is working"}))
		})

		app.Resource("/sounds", SoundsResource{})
		app.Resource("/users", UsersResource{})
		app.Resource("/threads", ThreadsResource{})
		app.Resource("/message_sounds", MessageSoundsResource{})
		app.Resource("/messages", MessagesResource{})

		
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}