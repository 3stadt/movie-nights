package imdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Config struct {
	ApiKey string
}

var searchURL = "https://imdb-api.com/%s/API/SearchMovie/%s/%s" // language, api key, term

func (c *Config) SearchMovie(lang, term string) (*SearchResult, error) {
	res, err := http.Get(fmt.Sprintf(searchURL, lang, c.ApiKey, term))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	var sRes SearchResult
	err = json.Unmarshal(body, &sRes)
	return &sRes, err
}
