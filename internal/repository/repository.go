package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/outstaff-task/internal/domain/models"
	"github.com/kartochnik010/outstaff-task/internal/repository/postgres"
	"github.com/sirupsen/logrus"
)

type Music interface {
	StoreMusic(ctx context.Context, music models.Music) (id uint64, err error)
	GetMusic(ctx context.Context, searchMeta *models.SearchMetadata) ([]models.Music, error)
	DeleteMusicByID(ctx context.Context, id uint64) error
	UpdateMusicByID(ctx context.Context, m models.Music) error
}

type Repository struct {
	Music
}

func NewRepository(db *pgxpool.Pool, l *logrus.Logger) Repository {
	return Repository{
		Music: postgres.NewMusicRepo(db),
	}
}
