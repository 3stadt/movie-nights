package imdb

type Genre struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Movie struct {
	MovieID       string  `json:"id,omitempty"`
	Image         string  `json:"image,omitempty"`
	Title         string  `json:"title,omitempty"`
	Plot          string  `json:"plot,omitempty"`
	PlotLocal     string  `json:"plotLocal,omitempty"`
	ReleaseDate   string  `json:"releaseDate,omitempty"`
	Runtime       string  `json:"runtimeStr,omitempty"`
	Genres        []Genre `json:"genreList,omitempty"`
	ContentRating string  `json:"contentRating,omitempty"`
	ErrorMessage  string  `json:"errorMessage,omitempty"`
}
