package entity

import "time"

type Movie struct {
	ID               uint64     `json:"id"`
	OriginalTitle    string     `json:"original_title"`
	OriginalLanguage string     `json:"original_language"`
	Overview         *string    `json:"overview"`
	PosterPath       *string    `json:"poster_path"`
	BackdropPath     *string    `json:"backdrop_path"`
	Adult            bool       `json:"adult"`
	ReleaseDate      *time.Time `json:"release_date"`
	Budget           *uint64    `json:"budget"`
	Revenue          *int64     `json:"revenue"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
