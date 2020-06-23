package runtime

import (
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"google.golang.org/grpc"
)

type (
	// MiddlewareFunc type definition
	MiddlewareFunc func(w http.ResponseWriter, r *http.Request) error

	// Custom error handler which is called on graphql result has an error
	GraphqlErrorHandler func(errs gqlerrors.FormattedErrors)
)

type GraphqlHandler interface {
	CreateConnection(context.Context) (*grpc.ClientConn, func(), error)
	GetMutations(*grpc.ClientConn) graphql.Fields
	GetQueries(*grpc.ClientConn) graphql.Fields
}

// ServeMux is struct can execute graphql request via incoming HTTP request.
// This is inspired from grpc-gateway implementation, thanks!
type ServeMux struct {
	middlewares  []MiddlewareFunc
	ErrorHandler GraphqlErrorHandler

	handlers []GraphqlHandler
}

// NewServeMux creates ServeMux pointer
func NewServeMux(ms ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		middlewares: ms,
		handlers:    make([]GraphqlHandler, 0),
	}
}

// AddHandler registers graphql handler which is built via plugin
func (s *ServeMux) AddHandler(h GraphqlHandler) error {
	if err := s.validateHandler(h); err != nil {
		return err
	}
	s.handlers = append(s.handlers, h)
	return nil
}

// Validate handler definition
func (s *ServeMux) validateHandler(h GraphqlHandler) error {
	queries := h.GetQueries(nil)
	mutations := h.GetMutations(nil)

	// If handler doesn't have any definitions, pass
	if len(queries) == 0 && len(mutations) == 0 {
		return nil
	}

	schemaConfig := graphql.SchemaConfig{}
	if len(queries) > 0 {
		schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		})
	}
	if len(mutations) > 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		})
	}

	// Try to generate Schema and check error
	if _, err := graphql.NewSchema(schemaConfig); err != nil {
		return fmt.Errorf("Schema validation error: %s", err)
	}
	return nil
}

// Use adds more middlwares
func (s *ServeMux) Use(ms ...MiddlewareFunc) *ServeMux {
	s.middlewares = append(s.middlewares, ms...)
	return s
}

// ServeHTTP implements http.Handler
func (s *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.middlewares {
		if err := m(w, r); err != nil {
			http.Error(w, "middleware error occurred: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ctx := r.Context()

	queries := graphql.Fields{}
	mutations := graphql.Fields{}
	for _, h := range s.handlers {
		c, closer, err := h.CreateConnection(ctx)
		if err != nil {
			http.Error(w, "failed to create grpc connection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer closer()

		for k, v := range h.GetQueries(c) {
			queries[k] = v
		}
		for k, v := range h.GetMutations(c) {
			mutations[k] = v
		}
	}

	schemaConfig := graphql.SchemaConfig{}
	if len(queries) > 0 {
		schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: queries,
		})
	}
	if len(mutations) > 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutations,
		})
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		http.Error(w, "failed to build schema: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := parseRequest(r)
	if err != nil {
		http.Error(w, "failed to parse request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		Context:        ctx,
	})

	if len(result.Errors) > 0 {
		if s.ErrorHandler != nil {
			s.ErrorHandler(result.Errors)
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "failed to marshal response JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(out)))
	w.WriteHeader(http.StatusOK)
	w.Write(out) // nolint: errcheck
}
