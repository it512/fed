package http

import (
	"net/http"

	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
)

func NewGraphqlHTTPHandler(schema *graphql.Schema, engine *engine.ExecutionEngine, enableART bool) http.Handler {
	return &GraphQLHTTPRequestHandler{
		schema:    schema,
		engine:    engine,
		enableART: enableART,
	}
}

type GraphQLHTTPRequestHandler struct {
	engine    *engine.ExecutionEngine
	schema    *graphql.Schema
	enableART bool
}

func (g *GraphQLHTTPRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.handleHTTP(w, r)
}
