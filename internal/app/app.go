package app

import (
	"flag"
	"fmt"
	"net/http"
	"new-shortner/internal/config"
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
	flag.Parse()

	urlsRepo := inmemory.NewUrls(lg)
	urlsService := service.NewUrls(urlsRepo)
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
