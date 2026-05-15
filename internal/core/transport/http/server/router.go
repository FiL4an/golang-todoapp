package core_http_server

import (
	"fmt"
	"net/http"

	core_http_midleware "github.com/FiL4an/golang-todoapp/internal/core/transport/http/midleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middlware  []core_http_midleware.Middleware
}

func NewAPIVersionRouter(apiVersion ApiVersion, middleware ...core_http_midleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: ApiVersion1,
		middlware:  middleware,
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_midleware.ChainMiddleware(r, r.middlware...)

}

func (r *APIVersionRouter) RegisterRouters(routers ...Route) {
	for _, route := range routers {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}
