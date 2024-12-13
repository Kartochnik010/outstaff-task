package service

import (
	"net/http"

	"github.com/kartochnik010/outstaff-task/internal/config"
	"github.com/kartochnik010/outstaff-task/internal/repository"
)

type Service struct {
	Music *MusicService
}

func NewService(repo repository.Repository, cfg *config.Config, c *http.Client) *Service {
	return &Service{
		Music: NewMusicService(c, cfg, repo),
	}
}
