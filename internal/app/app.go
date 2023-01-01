package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"new-shortner/internal/config"
	"new-shortner/internal/repository/file"
	"new-shortner/internal/repository/inmemory"
	"new-shortner/internal/service"
	"new-shortner/internal/transport/rest"
	"new-shortner/pkg/logger"

	"go.uber.org/zap"
)

type App struct {
	HTTPServer *http.Server
	logger     *zap.Logger
}

func New(cfg config.Config) (*App, error) {
	lg, err := logger.New(true)
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url for short urls")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file for save/load urls")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "dsn for database")
	flag.Parse()

	var repo service.URLRepository
	if cfg.FileStoragePath != "" {
		if repo, err = file.New(cfg.FileStoragePath, lg); err != nil {
			log.Fatal(err)
		}
	} else {
		repo = inmemory.NewURLs(lg)
	}
	urlsService := service.NewURLs(repo)
	handler := rest.NewHandler(urlsService, cfg)

	srv := &http.Server{
		Handler: handler.InitRouter(),
		Addr:    fmt.Sprint(cfg.ServerAddress),
	}

	return &App{
		HTTPServer: srv,
		logger:     lg,
	}, nil
}

func (app *App) Run() error {
	app.logger.Info("server started", zap.String("addr", app.HTTPServer.Addr))
	return app.HTTPServer.ListenAndServe()
}
