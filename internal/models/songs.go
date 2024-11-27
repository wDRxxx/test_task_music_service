package models

import (
	"time"
)

type Song struct {
	ID          int       `json:"id"`
	Song        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Text        string    `json:"text,omitempty"`
	Link        string    `json:"link,omitempty"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MusicInfo struct {
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Text        string    `json:"text,omitempty"`
	Link        string    `json:"link,omitempty"`
}
