package httpServer

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/wDRxxx/test-task/internal/api"
	"github.com/wDRxxx/test-task/internal/models"
	"github.com/wDRxxx/test-task/internal/service"
	"github.com/wDRxxx/test-task/internal/utils"
)

// @Summary      Получение списка песен
// @Description  Возвращает список всех песен по заданной странице и фильтру
// @Accept       json
// @Produce      json
// @Param		 page query int true "Страница"
// @Param		 group query string false "Группа-исполнитель"
// @Param		 song query string false "Название песни"
// @Success      200 {array} models.Song
// @Failure 	 400 {object} models.DefaultResponse
// @Router       /songs [get]
func (s *server) songs(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("page")
	page, err := strconv.Atoi(p)
	if err != nil || page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	songs, err := s.apiService.Songs(r.Context(), page, group, song)
	if err != nil {
		slog.Error("Error getting events", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(songs, w)
}

// @Summary      Возвращает песню
// @Description  Возвращает песню по ее ID
// @Accept       json
// @Produce      json
// @Success      200 {object} models.Song
// @Failure 	 400 {object} models.DefaultResponse
// @Param 		 id path string true "ID песни"
// @Router       /songs/{id} [get]
func (s *server) song(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strID)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	song, err := s.apiService.Song(r.Context(), id)
	if err != nil {
		slog.Error("Error getting song", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(song, w)
}

// @Summary      Создает новую песню
// @Description  Создает новую песню с заданными полями.  Обязательно проверяйте и изменяйте release_data для грамотной работы swagger'a
// @Accept       json
// @Produce      json
// @Success      200
// @Failure 	 400 {object} models.DefaultResponse
// @Param        song body models.Song true "Данные о песне"
// @Router       /songs [post]
func (s *server) createSong(w http.ResponseWriter, r *http.Request) {
	var song *models.Song
	err := utils.ReadJSON(r.Body, &song)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.apiService.CreateSong(r.Context(), song)
	if err != nil {
		slog.Error("Error creating song", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
}

// @Summary      Удаляет песню
// @Description  Создает новую песню с заданными полями
// @Produce      json
// @Success      200
// @Failure 	 400 {object} models.DefaultResponse
// @Param 		 id path string true "ID песни"
// @Router       /songs/{id} [delete]
func (s *server) deleteSong(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strID)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.apiService.DeleteSong(r.Context(), id)
	if err != nil {
		slog.Error("Error deleting song", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
}

// @Summary      Обновляет данные песни
// @Description  Обновляет данные песни.  Обязательно проверяйте и изменяйте release_data для грамотной работы swagger'a
// @Accept       json
// @Produce      json
// @Success      200
// @Failure 	 400 {object} models.DefaultResponse
// @Param 		 id path string true "ID песни"
// @Param        song body models.Song true "Данные о песне"
// @Router       /songs/{id} [patch]
func (s *server) updateSong(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strID)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var song *models.Song
	err = utils.ReadJSON(r.Body, &song)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteJSONError(err, w)
		return
	}
	song.ID = id

	err = s.apiService.UpdateSong(r.Context(), song)
	if err != nil {
		slog.Error("Error updating song", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}
}

// @Summary      Возвращает куплет песни
// @Description  Возвращает заданный куплет песни
// @Produce      json
// @Success      200 {object} models.DefaultResponse
// @Failure 	 400 {object} models.DefaultResponse
// @Param		 verse query string true "Номер куплета"
// @Param 		 id path string true "ID песни"
// @Router       /songs/{id}/text [get]
func (s *server) songVerse(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(strID)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strVerse := r.URL.Query().Get("verse")
	verseNum, err := strconv.Atoi(strVerse)
	if err != nil || verseNum < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	verse, err := s.apiService.SongVerse(r.Context(), id, verseNum-1)
	if err != nil {
		if errors.Is(err, service.ErrWrongVerse) {
			utils.WriteJSONError(err, w)
			return
		}

		slog.Error("Error getting song", slog.Any("error", err))
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.DefaultResponse{
		Error:   false,
		Message: verse,
	}, w)

}
