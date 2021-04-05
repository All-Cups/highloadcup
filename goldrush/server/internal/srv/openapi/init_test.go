package openapi_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/powerman/check"
	"github.com/powerman/structlog"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/smartystreets/goconvey/convey"

	"github.com/Djarvur/allcups-itrally-2020-task/api/openapi/client"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/app"
	"github.com/Djarvur/allcups-itrally-2020-task/internal/srv/openapi"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/def"
	"github.com/Djarvur/allcups-itrally-2020-task/pkg/netx"
)

func TestMain(m *testing.M) {
	def.Init()
	reg := prometheus.NewPedanticRegistry()
	app.InitMetrics(reg)
	openapi.InitMetrics(reg, "test")
	check.TestMain(m)
}

// Const shared by tests. Recommended naming scheme: <dataType><Variant>.
var (
	apiError402  = openapi.APIError(402, "bogus coin")
	apiError403  = openapi.APIError(403, "no such license")
	apiError404  = openapi.APIError(404, "no treasure")
	apiError429  = openapi.APIError(429, "too many requests")
	apiError500  = openapi.APIError(500, "internal error")
	apiError502  = openapi.APIError(502, "RPC failed")
	apiError503  = openapi.APIError(503, "service unavailable")
	apiError504  = openapi.APIError(504, "RPC timed out")
	apiError1000 = openapi.APIError(1000, "wrong coordinates")
	apiError1001 = openapi.APIError(1001, "wrong depth")
	apiError1002 = openapi.APIError(1002, "no more active licenses allowed")
	apiError1003 = openapi.APIError(1003, "treasure is not digged")
)

func testNewServer(t *check.C, cfg openapi.Config) (cleanup func(), c *client.HighLoadCup2020, url string, mockAppl *app.MockAppl, logc <-chan string) {
	cfg.Addr = netx.NewAddr("localhost", 0)

	t.Helper()
	ctrl := gomock.NewController(t)

	mockAppl = app.NewMockAppl(ctrl)
	mockAppl.EXPECT().Start(gomock.Any()).Return(nil).AnyTimes()

	server, err := openapi.NewServer(mockAppl, cfg)
	t.Must(t.Nil(err, "NewServer"))

	piper, pipew := io.Pipe()
	server.SetHandler(interceptLog(pipew, server.GetHandler()))
	logch := make(chan string, 64) // Keep some unread log messages.
	go func() {
		scanner := bufio.NewScanner(piper)
		for scanner.Scan() {
			select {
			default: // Do not hang test because of some unread log messages.
			case logch <- scanner.Text():
			}
		}
		close(logch)
	}()

	t.Must(t.Nil(server.Listen(), "server.Listen"))
	errc := make(chan error, 1)
	go func() { errc <- server.Serve() }()

	cleanup = func() {
		t.Helper()
		t.Nil(server.Shutdown(), "server.Shutdown")
		t.Nil(<-errc, "server.Serve")
		pipew.Close()
		ctrl.Finish()
	}

	ln, err := server.HTTPListener()
	t.Must(t.Nil(err, "server.HTTPListener"))
	c = client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes:  []string{"http"},
		Host:     ln.Addr().String(),
		BasePath: client.DefaultBasePath,
	})
	url = fmt.Sprintf("http://%s", ln.Addr().String())

	// Avoid race between server.Serve() and server.Shutdown().
	ctx, cancel := context.WithTimeout(context.Background(), def.TestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	t.Must(t.Nil(err))
	_, err = (&http.Client{}).Do(req)
	t.Must(t.Nil(err, "connect to service"))
	<-logch

	return cleanup, c, url, mockAppl, logch
}

func interceptLog(out io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := structlog.FromContext(r.Context(), nil)
		log.SetOutput(out)
		r = r.WithContext(structlog.NewContext(r.Context(), log))
		next.ServeHTTP(w, r)
	})
}
