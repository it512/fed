package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	gatewayHttp "github.com/it512/fed/http"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
)

type authedTransport struct {
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//req.Header.Set("Authorization", "bearer "+t.key)
	req.Header.Set("x-test", uuid.NewString())
	return t.wrapped.RoundTrip(req)
}

func startServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := &http.Client{
		Transport: &authedTransport{
			wrapped: http.DefaultTransport,
		},
	}

	datasourceWatcher := NewDatasourcePoller(httpClient, DatasourcePollerConfig{
		Services: []ServiceConfig{
			{Name: "reviews", URL: "http://82.157.165.187:30041/query"},
		},
		PollingInterval: 30 * time.Second,
	})

	enableART := false

	var gqlHandlerFactory HandlerFactoryFn = func(schema *graphql.Schema, engine *engine.ExecutionEngine) http.Handler {
		return gatewayHttp.NewGraphqlHTTPHandler(schema, engine, enableART)
	}

	gateway := NewGateway(ctx, gqlHandlerFactory, httpClient)

	datasourceWatcher.Register(gateway)
	go datasourceWatcher.Run(ctx)

	gateway.Ready()

	mux := http.NewServeMux()
	mux.Handle("/query", gateway)

	if err := http.ListenAndServe(":10009", mux); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startServer()
}
