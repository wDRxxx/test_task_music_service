package repository

import (
	"context"

	"github.com/wDRxxx/test-task/internal/models"
)

type Repository interface {
	Songs(ctx context.Context, page int, group string, song string) ([]*models.Song, error)
	Song(ctx context.Context, id int) (*models.Song, error)
	InsertSong(ctx context.Context, song *models.Song) error
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, song *models.Song) error
}
