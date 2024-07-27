package main

import (
	"context"
	"github.com/spf13/viper"
	"github.com/viktor8881/service-utilities/db"
	"github.com/viktor8881/service-utilities/http/server"
	"go.uber.org/zap"
	"log"
	"net/http"
	userDomain "simpleserver/domain/user"
	generated "simpleserver/generated/http/server"
	"simpleserver/inner/user"
	"simpleserver/middleware"
	"time"
)

func RegisterRoutes(ctx context.Context, mux *http.ServeMux, logger *zap.Logger,
) []func() {
	tr := server.NewTransport(mux)
	dbConfig := db.DatabaseConfig{
		DSN:                viper.GetString("database.dsn"),
		DBType:             viper.GetString("database.type"),
		SetMaxOpenConns:    viper.GetInt("database.max_open_conns"),
		SetMaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		SetConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime") * time.Second,
	}

	mongoDB, cFunction, err := db.NewMongoDb(ctx, dbConfig, logger)
	if err != nil {
		log.Fatal(err)
	}

	repository := userDomain.NewRepository(mongoDB)

	userService := user.NewService(repository)

	generated.ListUserEndpoint(
		tr,
		userService.ListUser,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	generated.ListUserByEmailEndpoint(
		tr,
		userService.ListUserByEmail,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	generated.GetUserEndpoint(
		tr,
		userService.GetUser,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	generated.CreateUserEndpoint(
		tr,
		userService.CreateUser,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	generated.UpdateUserEndpoint(
		tr,
		userService.UpdateUser,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	generated.DeleteUserEndpoint(
		tr,
		userService.DeleteUser,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
		server.LoggerMiddleware(logger),
	)

	return []func(){
		cFunction,
	}
}
