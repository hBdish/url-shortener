package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
	"url-shortener/tools/logger"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting", slog.String("env", cfg.Env))
	log.Debug("debug log")

	storage, err := postgres.New(cfg.Db)

	if err != nil {
		log.Error("failed to init db connection", logger.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveURL("https://yandex.ru", "ya")

	if err != nil {
		log.Error("failed to save in db", logger.Err(err))
		os.Exit(1)
	}

	log.Info("saved url", slog.Int64("id", id))

	url, err := storage.GetURL("google")

	log.Info("get url", slog.String("url", url))
	_ = storage
}

const (
	envLocal = "dev"
	envStage = "stage"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envStage:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
