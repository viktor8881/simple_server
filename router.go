package main

import (
	"context"
	"github.com/spf13/viper"
	"github.com/viktor8881/service-utilities/db"
	"github.com/viktor8881/service-utilities/http/server"
	"github.com/viktor8881/service-utilities/tbot"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"net/http"
	userDomain "simpleserver/domain/user"
	servergenerated "simpleserver/generated/http/server"
	tbotgenerated "simpleserver/generated/tbot"
	"simpleserver/inner/user"
	"simpleserver/middleware"
	"time"
)

func RegisterRoutes(ctx context.Context, mux *http.ServeMux, logger *zap.Logger) []func() {
	tr := server.NewTransport(mux)

	mongoDB, cFunction, err := newMongoDb(ctx, logger)
	if err != nil {
		logger.Fatal("failed create mongoDB", zap.Error(err))
	}

	telebot, err := newTbot(ctx, logger)
	if err != nil {
		logger.Fatal("failed create tbot", zap.Error(err))
	}

	repository := userDomain.NewRepository(mongoDB)

	userLogicService := user.NewService(repository)

	servergenerated.ListUserEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.ListUser,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	servergenerated.ListUserByEmailEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.ListUserByEmail,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	servergenerated.GetUserEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.GetUser,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	servergenerated.CreateUserEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.CreateUser,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	servergenerated.UpdateUserEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.UpdateUser,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	servergenerated.DeleteUserEndpoint(
		tr,
		server.DecodeRequest,
		userLogicService.DeleteUser,
		server.EncodeResponse,
		server.ErrorHandler,
		logger,
		middleware.AuthMiddleware(logger, viper.GetString("server.api_token")),
	)

	tbotgenerated.ListUserEndpoint(
		ctx,
		telebot,
		tbot.DecodePayload,
		userLogicService.ListUser,
		user.OutListUserHelper(logger),
		tbot.ErrorHandler,
		logger,
	)

	tbotgenerated.ListUserByEmailEndpoint(
		ctx,
		telebot,
		user.InListUserByEmail,
		userLogicService.ListUserByEmail,
		user.OutListUserHelper(logger),
		tbot.ErrorHandler,
		logger,
	)

	tbotgenerated.GetUserEndpoint(
		ctx,
		telebot,
		tbot.DecodePayload,
		userLogicService.GetUser,
		tbot.EncodeResponse,
		tbot.ErrorHandler,
		logger,
	)

	go telebot.Start()

	return []func(){
		cFunction,
		telebot.Stop,
	}
}

func newMongoDb(ctx context.Context, logger *zap.Logger) (*db.MongoDB, func(), error) {
	dbConfig := db.DatabaseConfig{
		DSN:                viper.GetString("database.dsn"),
		DBType:             viper.GetString("database.type"),
		SetMaxOpenConns:    viper.GetInt("database.max_open_conns"),
		SetMaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		SetConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime") * time.Second,
	}

	return db.NewMongoDb(ctx, dbConfig, logger)
}

func newTbot(ctx context.Context, logger *zap.Logger) (*tbot.Bot, error) {
	pref := tele.Settings{
		Token:  viper.GetString("tbot.token"),
		Poller: &tele.LongPoller{Timeout: viper.GetDuration("tbot.timeout") * time.Second},
	}

	return tbot.NewBot(pref)
}
