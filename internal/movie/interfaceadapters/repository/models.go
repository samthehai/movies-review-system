package repository

import "time"

type Movie struct {
	ID               uint64     `json:"id" db:"id"`
	OriginalTitle    string     `json:"original_title" db:"original_title"`
	OriginalLanguage string     `json:"original_language" db:"original_language"`
	Overview         *string    `json:"overview" db:"overview"`
	PosterPath       *string    `json:"poster_path" db:"poster_path"`
	BackdropPath     *string    `json:"backdrop_path" db:"backdrop_path"`
	Adult            bool       `json:"adult" db:"adult"`
	ReleaseDate      *time.Time `json:"release_date" db:"release_date"`
	Budget           *uint64    `json:"budget" db:"budget"`
	Revenue          *int64     `json:"revenue" db:"revenue"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
