package imdb

import "time"

type Genre struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Movie struct {
	MovieID       string         `json:"id,omitempty"`
	Image         string         `json:"image,omitempty"`
	Title         string         `json:"title,omitempty"`
	Plot          string         `json:"plotLocal,omitempty"`
	ReleaseDate   *time.Time     `json:"releaseDate,omitempty"`
	Runtime       *time.Duration `json:"runtimeMins,omitempty"`
	GenreList     []Genre        `json:"genreList,omitempty"`
	ContentRating string         `json:"contentRating,omitempty"`
}
