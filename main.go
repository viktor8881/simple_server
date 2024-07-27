package main

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"time"

	appHttp "github.com/viktor8881/service-utilities/http"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
