package httpServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/wDRxxx/test-task/docs"
)

func (s *server) setRoutes() {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/swagger/*", httpSwagger.WrapHandler)

	mux.Route("/songs", func(mux chi.Router) {
		mux.Get("/", s.songs)
		mux.Post("/", s.createSong)

		mux.Route("/{id}", func(mux chi.Router) {
			mux.Get("/text", s.songVerse)

			mux.Get("/", s.song)
			mux.Delete("/", s.deleteSong)
			mux.Patch("/", s.updateSong)
		})
	})

	s.mux = mux
}
