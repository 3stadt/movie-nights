package imdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var movieURL = "https://imdb-api.com/%s/API/Title/%s/%s" // language, api key, movie id

func (c *Config) MovieDetail(lang, movieID string) (*Movie, error) {
	res, err := http.Get(fmt.Sprintf(movieURL, lang, c.ApiKey, movieID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var m Movie
	err = json.Unmarshal(body, &m)
	return &m, err
}
