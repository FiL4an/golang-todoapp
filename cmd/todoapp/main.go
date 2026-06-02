package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/FiL4an/golang-todoapp/internal/core/config"
	core_logger "github.com/FiL4an/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/FiL4an/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_midleware "github.com/FiL4an/golang-todoapp/internal/core/transport/http/midleware"
	core_http_server "github.com/FiL4an/golang-todoapp/internal/core/transport/http/server"
	task_postgres_repository "github.com/FiL4an/golang-todoapp/internal/feature/tasks/repository/postgres"
	task_service "github.com/FiL4an/golang-todoapp/internal/feature/tasks/service"
	task_transport_http "github.com/FiL4an/golang-todoapp/internal/feature/tasks/transport/http"
	users_postgres_repository "github.com/FiL4an/golang-todoapp/internal/feature/users/repository/postgres"
	users_service "github.com/FiL4an/golang-todoapp/internal/feature/users/service"
	users_transport_http "github.com/FiL4an/golang-todoapp/internal/feature/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM)

	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init applications:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))
	logger.Debug("initazling postgres connection pool")

	pool, err := core_pgx_pool.NewConnectionPool(core_pgx_pool.NewConfigMust(), ctx)
	if err != nil {
		logger.Fatal("failed to init postgres pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature ", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUserRepository(pool)
	userService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(userService)

	logger.Debug("initializing", zap.String("feature", "tasks"))

	tasksRepository := task_postgres_repository.NewTasksRepository(pool)
	tasksService := task_service.NewTasksService(tasksRepository)
	taskTransportHTTP := task_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_midleware.RequestID(),
		core_http_midleware.Logger(logger),
		core_http_midleware.Trace(),
		core_http_midleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRouters(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRouters(taskTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error(" HTTP server run error", zap.Error(err))
	}
}
