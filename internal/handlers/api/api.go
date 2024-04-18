package api

import (
	"context"
	swagger "effectiveM-test-task/internal/api"
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct{
	Path string `env:"BASEPATH"` 
}

func SomeApi(ctx context.Context, log *slog.Logger) error{
	cfg := getBasePath()

	conf := swagger.NewConfiguration()
	conf.BasePath = cfg.Path

	apiClient := swagger.NewAPIClient(conf)
	_, resp, err := apiClient.DefaultApi.InfoGet(ctx, "X123XX150")
	if err!= nil {
		log.Error("API: %v", err)
		return nil
	}

	msg := fmt.Sprintf(" SomeAPI called: Status: %v", resp.Status) 
	log.Debug(msg)
	return nil
}


func getBasePath() Config{
	var cfg Config
	err := cleanenv.ReadConfig("config/api.env" ,&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}