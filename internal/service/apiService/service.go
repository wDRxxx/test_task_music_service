package apiService

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/wDRxxx/test-task/internal/models"
	"github.com/wDRxxx/test-task/internal/repository"
	"github.com/wDRxxx/test-task/internal/service"
	"github.com/wDRxxx/test-task/internal/utils"
)

type serv struct {
	repo repository.Repository

	client              http.Client
	musicInfoServiceURL string
}

func NewApiService(repo repository.Repository, musicInfoServiceURL string) service.ApiService {
	return &serv{
		repo:                repo,
		client:              http.Client{},
		musicInfoServiceURL: musicInfoServiceURL,
	}
}

func (s *serv) Songs(ctx context.Context, page int, group string, song string) ([]*models.Song, error) {
	songs, err := s.repo.Songs(ctx, page, group, song)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *serv) Song(ctx context.Context, id int) (*models.Song, error) {
	song, err := s.repo.Song(ctx, id)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *serv) SongVerse(ctx context.Context, id int, verse int) (string, error) {
	song, err := s.repo.Song(ctx, id)
	if err != nil {
		return "", err
	}

	verses := strings.Split(song.Text, "\n")
	if len(verses) <= verse {
		return "", service.ErrWrongVerse
	}

	return verses[verse], nil
}

func (s *serv) CreateSong(ctx context.Context, song *models.Song) error {
	url := fmt.Sprintf(
		"%s/info?group=%s&song=%s",
		s.musicInfoServiceURL, song.Group, song.Group,
	)
	resp, err := s.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var musicInfo models.MusicInfo
	err = utils.ReadJSON(resp.Body, &musicInfo)
	if err != nil {
		return err
	}

	song.ReleaseDate = musicInfo.ReleaseDate
	song.Text = musicInfo.Text
	song.Link = musicInfo.Link

	err = s.repo.InsertSong(ctx, song)
	if err != nil {
		return err
	}

	slog.Info("new song was created", slog.Any("song", song))

	return nil
}

func (s *serv) DeleteSong(ctx context.Context, id int) error {
	err := s.repo.DeleteSong(ctx, id)
	if err != nil {
		return err
	}

	slog.Info("song was deleted", slog.Int("id", id))

	return nil
}

func (s *serv) UpdateSong(ctx context.Context, song *models.Song) error {
	err := s.repo.UpdateSong(ctx, song)
	if err != nil {
		return err
	}

	slog.Info("song was updated", slog.Any("song", song))

	return nil
}
