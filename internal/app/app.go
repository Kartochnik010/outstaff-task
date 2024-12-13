package app

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/outstaff-task/internal/config"
	"github.com/kartochnik010/outstaff-task/internal/handler"
	"github.com/kartochnik010/outstaff-task/internal/repository"
	"github.com/kartochnik010/outstaff-task/internal/service"
	"github.com/sirupsen/logrus"
)

type App struct {
	Server *http.Server
	Logger *logrus.Logger
	Repo   repository.Repository
}

func NewApp(cfg *config.Config, db *pgxpool.Pool, l *logrus.Logger) *App {
	const op = "app.NewApp"

	repo := repository.NewRepository(db, l)

	service := service.NewService(repo, cfg, &http.Client{})

	h := handler.NewHandler(&http.Client{}, service, l, cfg)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: handler.Routes(h),
	}
	return &App{
		Server: s,
		Logger: l,
		Repo:   repo,
	}
}

func (a *App) Run() error {
	return a.Server.ListenAndServe()
}
