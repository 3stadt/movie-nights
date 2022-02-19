package imdb

type SearchResultItem struct {
	MovieID     string `json:"id,omitempty"`
	Type        string `json:"resultType,omitempty"`
	Image       string `json:"image,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type SearchResult struct {
	Results      []SearchResultItem `json:"results,omitempty"`
	Type         string             `json:"type,omitempty"`
	Expression   string             `json:"expression,omitempty"`
	ErrorMessage string             `json:"errorMessage,omitempty"`
}
