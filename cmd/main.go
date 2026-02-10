package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/broadcast80/ozon-task/config"
	"github.com/broadcast80/ozon-task/domain/link"
	app "github.com/broadcast80/ozon-task/internal/app"
	"github.com/broadcast80/ozon-task/internal/pkg/utils"
	inmemory "github.com/broadcast80/ozon-task/internal/repository/in_memory"
	"github.com/broadcast80/ozon-task/internal/repository/postgresql"
	"github.com/broadcast80/ozon-task/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		print("ENV_PATH required")
		return
	}

	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("error loading .env file from path: %s", envPath)
	}

	cfg := config.MustLoad()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("starting service")
	log.Debug("debug messages are enabled")

	ctx := context.TODO()

	repository := newRepository(ctx, *cfg, log)

	dataProvider := usecase.New(repository, log)

	service := link.NewShortener(dataProvider)

	router := http.NewServeMux()

	handlers := app.New(router, service, log)

	if err = handlers.MapHandlers(); err != nil {
		log.Error("failed to map handlers")
	}

	errs := make(chan error, 2)

	go func() {
		errs <- handlers.ListenAndServe(cfg.HTTPServer.Port)
	}()

	err = <-errs
	if err != nil {
		fmt.Printf("Werr %s", err.Error())
	}
}

func newRepository(ctx context.Context, cfg config.Config, log *slog.Logger) usecase.RepositoryInterface {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "inmemory"
	}

	switch storageType {

	case "postgres":
		postgreSQLClient, err := utils.NewClient(ctx, 5, cfg.PostgresConfig)
		if err != nil {
			log.Error("failed to init storage", "Error", err.Error())
			os.Exit(1)
		}
		repository := postgresql.New(postgreSQLClient)
		return repository

	case "inmemory":
		repository := inmemory.New(cfg.InMemoryConfig.Size)
		return repository

	default:
		log.Error("uknown STORAGE_TYPE", "STORAGE_TYPE", storageType)
		return nil
	}
}
