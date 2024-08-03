package main

import (
	"context"
	"github.com/spf13/viper"
	appHttp "github.com/viktor8881/service-utilities/http"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed create logger: %w", err)
	}
	defer logger.Sync()

	err = LoadConfig()
	if err != nil {
		log.Fatalf("failed load config: %w", err)
	}

	app := appHttp.NewApp(viper.GetString("server.host")+":"+viper.GetString("server.port"), logger)
	closeFuncs := RegisterRoutes(ctx, app.Mux, logger)
	for _, f := range closeFuncs {
		defer f()
	}
	app.Run()
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
