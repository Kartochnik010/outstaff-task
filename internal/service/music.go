package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kartochnik010/outstaff-task/internal/config"
	"github.com/kartochnik010/outstaff-task/internal/domain/models"
	"github.com/kartochnik010/outstaff-task/internal/pkg/logger"
	"github.com/kartochnik010/outstaff-task/internal/repository"
)

type MusicService struct {
	c    *http.Client
	cfg  *config.Config
	repo repository.Repository
}

func NewMusicService(client *http.Client, cfg *config.Config, repo repository.Repository) *MusicService {
	return &MusicService{
		c:    client,
		cfg:  cfg,
		repo: repo,
	}
}

func (s *MusicService) StoreMusic(ctx context.Context, music models.Music) (uint64, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicService.StoreMusic")

	//		fetch music info from music api
	url := s.cfg.MusicApiBaseUrl + fmt.Sprintf(`/info?group="%s"&&song="%s"`, music.Group, music.Song)
	resp, err := s.c.Get(url)
	if err != nil {
		log.WithError(err).Error("failed to fetch musci")
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to read body")
		return 0, err
	}
	if err = json.Unmarshal(body, &music); err != nil {
		log.WithError(err).Error("failed to unmarshal")
		return 0, err
	}
	// validate

	// save rich music info
	id, err := s.repo.Music.StoreMusic(ctx, music)
	if err != nil {
		log.WithError(err).Error("failed to store music")
		return 0, err
	}
	return id, nil
}

func (s *MusicService) GetMusic(ctx context.Context, searchMeta *models.SearchMetadata) ([]models.Music, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicService.GetMusic")

	music, err := s.repo.Music.GetMusic(ctx, searchMeta)
	if err != nil {
		log.WithError(err).Error("failed to get music")
		return nil, err
	}

	return music, nil
}

func (s *MusicService) DeleteMusic(ctx context.Context, id uint64) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicService.DeleteMusic")

	err := s.repo.Music.DeleteMusicByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("failed to delete music")
		return err
	}

	return nil
}
func (s *MusicService) UpdateMusicByID(ctx context.Context, m models.Music) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicService.DeleteMusic")

	err := s.repo.Music.UpdateMusicByID(ctx, m)
	if err != nil {
		log.WithError(err).Error("failed to delete music")
		return err
	}

	return nil
}
