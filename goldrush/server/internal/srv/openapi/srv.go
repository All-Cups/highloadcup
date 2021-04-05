// Package openapi implements OpenAPI server.
package openapi

import (
	"context"
	"fmt"
	prng "math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/powerman/structlog"
	"github.com/sebest/xff"
	"golang.org/x/time/rate"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi"
	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/restapi/op"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

type (
	// Ctx is a synonym for convenience.
	Ctx = context.Context
	// Log is a synonym for convenience.
	Log = *structlog.Logger
	// Config contains configuration for OpenAPI server.
	Config struct {
		Seed                 int64
		DisableAccessLog     bool
		Addr                 netx.Addr
		BasePath             string
		Pprof                bool
		OpCashPercentFail    int
		OpCashRate           int
		OpDigRate            int
		OpDigTimeout         time.Duration
		OpExploreAreaRate    int
		OpExploreAreaTimeout time.Duration
		OpGetBalanceRate     int
		OpIssueLicenseRate   int
		OpListLicensesRate   int
	}
	server struct {
		app               app.Appl
		cfg               Config
		muPRNG            sync.Mutex
		prng              *prng.Rand
		limitGetBalance   *rate.Limiter
		limitListLicenses *rate.Limiter
		limitIssueLicense *rate.Limiter
		limitExploreArea  *rate.Limiter
		limitDig          *rate.Limiter
		limitCash         *rate.Limiter
	}
)

// NewServer returns OpenAPI server configured to listen on the TCP network
// address cfg.Host:cfg.Port and handle requests on incoming connections.
func NewServer(appl app.Appl, cfg Config) (*restapi.Server, error) {
	if cfg.Seed == 0 {
		cfg.Seed = time.Now().UnixNano()
	}
	srv := &server{
		app:               appl,
		cfg:               cfg,
		prng:              prng.New(prng.NewSource(cfg.Seed)), //nolint:gosec // We need repeatable results.
		limitGetBalance:   rate.NewLimiter(rate.Limit(cfg.OpGetBalanceRate), cfg.OpGetBalanceRate),
		limitListLicenses: rate.NewLimiter(rate.Limit(cfg.OpListLicensesRate), cfg.OpListLicensesRate),
		limitIssueLicense: rate.NewLimiter(rate.Limit(cfg.OpIssueLicenseRate), cfg.OpIssueLicenseRate),
		limitExploreArea:  rate.NewLimiter(rate.Limit(cfg.OpExploreAreaRate), cfg.OpExploreAreaRate),
		limitDig:          rate.NewLimiter(rate.Limit(cfg.OpDigRate), cfg.OpDigRate),
		limitCash:         rate.NewLimiter(rate.Limit(cfg.OpCashRate), cfg.OpCashRate),
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, fmt.Errorf("load embedded swagger spec: %w", err)
	}
	if cfg.BasePath == "" {
		cfg.BasePath = swaggerSpec.BasePath()
	}
	swaggerSpec.Spec().BasePath = cfg.BasePath

	api := op.NewHighLoadCup2020API(swaggerSpec)
	api.Logger = structlog.New(structlog.KeyUnit, "swagger").Printf

	api.HealthCheckHandler = op.HealthCheckHandlerFunc(srv.HealthCheck)
	api.GetBalanceHandler = op.GetBalanceHandlerFunc(srv.GetBalance)
	api.ListLicensesHandler = op.ListLicensesHandlerFunc(srv.ListLicenses)
	api.IssueLicenseHandler = op.IssueLicenseHandlerFunc(srv.IssueLicense)
	api.ExploreAreaHandler = op.ExploreAreaHandlerFunc(srv.ExploreArea)
	api.DigHandler = op.DigHandlerFunc(srv.Dig)
	api.CashHandler = op.CashHandlerFunc(srv.Cash)

	server := restapi.NewServer(api)
	server.CleanupTimeout = 10 * time.Second //nolint:gomnd // Const.
	server.GracefulTimeout = 5 * time.Second //nolint:gomnd // Const.
	server.MaxHeaderSize = 1000000           //nolint:gomnd // Const.
	server.KeepAlive = 3 * time.Minute       //nolint:gomnd // Const.
	server.ReadTimeout = 30 * time.Second    //nolint:gomnd // Const.
	server.WriteTimeout = 30 * time.Second   //nolint:gomnd // Const.
	server.TLSKeepAlive = server.KeepAlive
	server.TLSReadTimeout = server.ReadTimeout
	server.TLSWriteTimeout = server.WriteTimeout
	server.Host = cfg.Addr.Host()
	server.Port = cfg.Addr.Port()

	// The middleware executes before anything.
	api.UseSwaggerUI()
	globalMiddlewares := func(handler http.Handler) http.Handler {
		xffmw, _ := xff.Default()
		logger := makeLogger(cfg.BasePath)
		accesslog := makeAccessLog(cfg.BasePath, cfg.DisableAccessLog)
		optPprof := func(next http.Handler) http.Handler { return next }
		if cfg.Pprof {
			optPprof = pprof
		}
		return noCache(xffmw.Handler(logger(recovery(accesslog(optPprof(
			middleware.Spec(cfg.BasePath, restapi.FlatSwaggerJSON,
				cors(handler))))))))
	}
	// The middleware executes after serving /swagger.json and routing,
	// but before authentication, binding and validation.
	middlewares := func(handler http.Handler) http.Handler {
		appStart := makeAppStart(cfg.BasePath, srv.app)
		return appStart(handler)
	}
	server.SetHandler(globalMiddlewares(api.Serve(middlewares)))

	log := structlog.New()
	log.Info("OpenAPI protocol", "version", swaggerSpec.Spec().Info.Version)
	return server, nil
}

func (srv *server) inPercent(p int) bool {
	const percent100 = 100
	srv.muPRNG.Lock()
	defer srv.muPRNG.Unlock()
	return srv.prng.Intn(percent100) < p
}

func fromRequest(r *http.Request) (Ctx, Log) {
	ctx := r.Context()
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	ctx = def.NewContextWithRemoteIP(ctx, remoteIP)
	log := structlog.FromContext(ctx, nil)
	return ctx, log
}
