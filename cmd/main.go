package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/broadcast80/ozon-task/internal/application"
	"github.com/broadcast80/ozon-task/internal/config"
	infrastracture "github.com/broadcast80/ozon-task/internal/infrastructure/http"
	"github.com/broadcast80/ozon-task/internal/infrastructure/repository"
	"github.com/broadcast80/ozon-task/internal/pkg/db/postgresql"
	"github.com/joho/godotenv"
)

func main() {

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		print("posos")
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

	postgreSQLClient, err := postgresql.NewClient(ctx, 5, cfg.PostgresConfig)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	repository := repository.New(postgreSQLClient)

	service := application.New(repository)

	router := http.NewServeMux()

	handlers := infrastracture.New(cfg, router, service)

	if err = handlers.MapHandlers(); err != nil {
		log.Error("failed to map handlers")
	}

	errs := make(chan error, 2)

	go func() {
		errs <- handlers.ListenAndServe(cfg.HTTPServer)
	}()

	err = <-errs
	if err != nil {
		fmt.Printf("Werr %s", err.Error())
	}
}
