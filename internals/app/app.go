package app

import (
	"bannerService/internals/config"
	"bannerService/internals/handlers"
	"bannerService/internals/service"
	"bannerService/internals/storage"
	"bannerService/internals/storage/migrator"
	"bannerService/internals/storage/postgres"
	"log/slog"
	"net/http"
)

type App struct {
	Cfg     *config.Config
	Logger  *slog.Logger
	Storage storage.Storage
	Service *service.Service
	Server  *handlers.Handler
}

func New(cfg *config.Config, logger *slog.Logger) *App {

	dbConnect := postgres.GetDBConnect(cfg, logger)

	st := postgres.NewPostgresStorage(dbConnect)

	svc := service.NewService(st, logger)

	srv := handlers.NewHandler(svc, logger, cfg)

	return &App{
		Cfg:     cfg,
		Logger:  logger,
		Storage: st,
		Service: svc,
		Server:  srv,
	}
}

func (a *App) MustStart() {

	migrator.MustRunMigrations(a.Cfg, a.Logger)

	mux := http.NewServeMux()
	a.Server.RegisterRoutes(mux)

	a.Logger.Info("starting HTTP server", slog.String("addr", a.Cfg.App.ServerAddr))
	err := http.ListenAndServe(a.Cfg.App.ServerAddr, mux)
	if err != nil {
		a.Logger.Error("failed to start HTTP server", slog.String("error", err.Error()))
	}

}

func (a *App) Close() {
	a.Storage.Close()
}
