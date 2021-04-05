// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi/op"
)


func configureFlags(api *op.HighLoadCup2020API) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *op.HighLoadCup2020API) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.CashHandler == nil {
		api.CashHandler = op.CashHandlerFunc(func(params op.CashParams) op.CashResponder {
			return op.CashNotImplemented()
		})
	}
	if api.DigHandler == nil {
		api.DigHandler = op.DigHandlerFunc(func(params op.DigParams) op.DigResponder {
			return op.DigNotImplemented()
		})
	}
	if api.ExploreAreaHandler == nil {
		api.ExploreAreaHandler = op.ExploreAreaHandlerFunc(func(params op.ExploreAreaParams) op.ExploreAreaResponder {
			return op.ExploreAreaNotImplemented()
		})
	}
	if api.GetBalanceHandler == nil {
		api.GetBalanceHandler = op.GetBalanceHandlerFunc(func(params op.GetBalanceParams) op.GetBalanceResponder {
			return op.GetBalanceNotImplemented()
		})
	}
	if api.HealthCheckHandler == nil {
		api.HealthCheckHandler = op.HealthCheckHandlerFunc(func(params op.HealthCheckParams) op.HealthCheckResponder {
			return op.HealthCheckNotImplemented()
		})
	}
	if api.IssueLicenseHandler == nil {
		api.IssueLicenseHandler = op.IssueLicenseHandlerFunc(func(params op.IssueLicenseParams) op.IssueLicenseResponder {
			return op.IssueLicenseNotImplemented()
		})
	}
	if api.ListLicensesHandler == nil {
		api.ListLicensesHandler = op.ListLicensesHandlerFunc(func(params op.ListLicensesParams) op.ListLicensesResponder {
			return op.ListLicensesNotImplemented()
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
