package core_http_server

import (
	"net/http"

	core_http_midleware "github.com/FiL4an/golang-todoapp/internal/core/transport/http/midleware"
)

type Route struct {
	Method    string
	Path      string
	Handler   http.HandlerFunc
	Midleware []core_http_midleware.Middleware
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_midleware.ChainMiddleware(r.Handler, r.Midleware...)
}
