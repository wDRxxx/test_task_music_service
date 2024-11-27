package postgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/test-task/internal/models"
	"github.com/wDRxxx/test-task/internal/repository"
)

type repo struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewPostgresRepo(db *pgxpool.Pool, timeout time.Duration) repository.Repository {
	return &repo{
		db:      db,
		timeout: timeout,
	}
}

func (r *repo) Songs(ctx context.Context, page int, group string, song string) ([]*models.Song, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(`id, "group", song, release_date, text, link, created_at, updated_at`).
		From("songs").
		Offset(uint64((page - 1) * 10)).
		Limit(10).
		PlaceholderFormat(sq.Dollar)

	if group != "" {
		builder = builder.Where(sq.Eq{`"group"`: group})
	}
	if song != "" {
		builder = builder.Where(sq.Eq{"song": song})
	}

	sql, args, err := builder.ToSql()

	slog.Debug("running sql", slog.String("sql", sql), slog.Any("args", args))

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}
	defer rows.Close()

	songs := make([]*models.Song, 0, 10)
	for rows.Next() {
		var song models.Song
		err = rows.Scan(
			&song.ID,
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		songs = append(songs, &song)
	}

	return songs, nil
}

func (r *repo) Song(ctx context.Context, id int) (*models.Song, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(`id, "group", song, release_date, text, link, created_at, updated_at`).
		From("songs").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()

	slog.Debug("running sql", slog.String("sql", sql), slog.Any("args", args))

	var song models.Song
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&song.ID,
		&song.Group,
		&song.Song,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &song, nil
}

func (r *repo) InsertSong(ctx context.Context, song *models.Song) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Insert("songs").
		Columns("song", `"group"`, "release_date", "text", "link").
		Values(song.Song, song.Group, song.ReleaseDate, song.Text, song.Link).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	slog.Debug("running sql", slog.String("sql", sql), slog.Any("args", args))

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteSong(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Delete("songs").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()

	slog.Debug("running sql", slog.String("sql", sql), slog.Any("args", args))

	_, err = r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateSong(ctx context.Context, song *models.Song) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Update("songs").
		Where(sq.Eq{"id": song.ID}).
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now())

	if song.Song != "" {
		builder = builder.Set("song", song.Song)
	}
	if song.Group != "" {
		builder = builder.Set(`"group"`, song.Group)
	}
	if !song.ReleaseDate.IsZero() {
		builder = builder.Set("release_date", song.ReleaseDate)
	}
	if song.Text != "" {
		builder = builder.Set("text", song.Text)
	}
	if song.Link != "" {
		builder = builder.Set("link", song.Link)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	slog.Debug("running sql", slog.String("sql", sql), slog.Any("args", args))

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
