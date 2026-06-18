package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/FiL4an/golang-todoapp/docs"
	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_http_midleware "github.com/FiL4an/golang-todoapp/internal/core/transport/http/midleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_midleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_midleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		h.mux.Handle(pattern, route.WithMiddleware())
	}
}

func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (h *HTTPServer) RegisterSwagger() {
	h.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)
	h.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

// func (h *HTTPServer) RegisterPublic() {
// 	h.mux.Handle("/", http.FileServer(http.Dir("./public")))
// }

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_midleware.ChainMiddleware(h.mux, h.middleware...)
	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("start HTTP server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {

	case err := <-ch:
		if err != nil {
			return fmt.Errorf("Listen and serve HTTP : %w", err)
		}

	case <-ctx.Done():
		h.log.Warn("Shutdown HTTP server ")

		shutdownctx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownctx); err != nil {
			_ = server.Close()
			fmt.Errorf("Failed shutdown HTTP server: %w", err)
		}
		h.log.Warn("HTTP server stoped")
	}
	return nil
}
