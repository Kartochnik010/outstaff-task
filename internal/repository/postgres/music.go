package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kartochnik010/outstaff-task/internal/domain"
	"github.com/kartochnik010/outstaff-task/internal/domain/models"
	"github.com/kartochnik010/outstaff-task/internal/pkg/logger"
)

type MusicRepo struct {
	db *pgxpool.Pool
}

func NewMusicRepo(db *pgxpool.Pool) *MusicRepo {
	return &MusicRepo{
		db: db,
	}
}

func (r *MusicRepo) StoreMusic(ctx context.Context, m models.Music) (uint64, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicRepo.StoreMusic")

	query := `
		INSERT INTO music (group_name, song, link, text, release_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id uint64
	err := r.db.QueryRow(ctx, query, m.Group, m.Song, m.Link, m.Text, m.ReleaseDate).Scan(&id)
	if err != nil {
		log.WithError(err).Error("failed to save item")
		return 0, domain.ErrInternal
	}
	m.ID = id
	log.Debugf("saved music: %+v", m)
	return id, nil
}

func (c *MusicRepo) GetMusic(ctx context.Context, meta *models.SearchMetadata) ([]models.Music, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicRepo.GetMusic")
	log.Debugf("getting music with meta: %+v\n", meta)

	query, args := buildArgs(meta)
	offset := meta.Page*meta.Limit - meta.Limit
	args = append([]any{meta.Limit, offset}, args...)

	//get by id
	if meta.ID != 0 {
		query = `
		SELECT 
			id, group_name, song, link, text, release_date
		FROM music
		WHERE id = $1
	`
		args = []interface{}{meta.ID}
	}

	log.Debugf("query: %s, args: %v", query, args)
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		log.WithError(err).Error("failed to get music")
		return nil, err
	}
	defer rows.Close()

	res := []models.Music{}
	for rows.Next() {
		m := models.Music{}
		err := rows.Scan(&m.ID, &m.Group, &m.Song, &m.Link, &m.Text, &m.ReleaseDate)
		if err != nil {
			log.WithError(err).Error("failed to scan row")
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}

func buildArgs(meta *models.SearchMetadata) (query string, args []interface{}) {
	query = `SELECT id, group_name, song, link, text, release_date FROM music `
	where := ` $$ = '' OR (lower(%s) LIKE '%%' || lower($$) || '%%') `
	wheres := []string{}
	if meta.ID != 0 {
		wheres = append(wheres, fmt.Sprintf(where, "id"))
		args = append(args, meta.ID)
	}
	if meta.Group != "" {
		wheres = append(wheres, fmt.Sprintf(where, "group_name"))
		args = append(args, meta.Group)
	}
	if meta.Song != "" {
		wheres = append(wheres, fmt.Sprintf(where, "song"))
		args = append(args, meta.Song)
	}
	if meta.Link != "" {
		wheres = append(wheres, fmt.Sprintf(where, "link"))
		args = append(args, meta.Link)
	}
	if meta.Text != "" {
		wheres = append(wheres, fmt.Sprintf(where, "text"))
		args = append(args, meta.Text)
	}
	if meta.ReleaseDate != nil && !meta.ReleaseDate.IsZero() {
		wheres = append(wheres, fmt.Sprintf(where, "release_date"))
		args = append(args, meta.ReleaseDate)
	}
	if len(wheres) > 0 {
		query += "WHERE "
		for i, v := range wheres {
			query += strings.ReplaceAll(v, "$$", fmt.Sprintf("$%v", i+1+2)) + "AND"
		}
		query = strings.TrimSuffix(query, "AND")
	}
	return query + " LIMIT $1 OFFSET $2", args
}

func (r *MusicRepo) DeleteMusicByID(ctx context.Context, id uint64) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicRepo.GetMusicByID")

	query := `
		DELETE FROM music
		WHERE id = $1
		`

	tags, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.WithError(err).Error("failed to get music by id")
		return err
	}

	if tags.RowsAffected() != 1 {
		return domain.ErrMusicNotFound
	}

	log.Debugf("deleting music with id: %v", id)
	return nil
}

func (r *MusicRepo) UpdateMusicByID(ctx context.Context, m models.Music) error {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "MusicRepo.UpdateMusicByID")

	query := `
		UPDATE music
		SET group_name = $1, song = $2, link = $3, text = $4, release_date = $5
		WHERE id = $6
		
		`

	tags, err := r.db.Exec(ctx, query, m.Group, m.Song, m.Link, m.Text, m.ReleaseDate, m.ID)
	if err != nil {
		log.WithError(err).Error("failed to get music by id")
		return err
	}

	found := tags.RowsAffected() == 1
	if !found {
		return domain.ErrMusicNotFound
	}

	log.Debugf("updating music with id: %v", m.ID)
	return nil
}
