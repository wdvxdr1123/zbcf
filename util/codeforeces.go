package util

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type H = map[string]string

type API struct {
	url string
}

func NewCodeforcesAPI(url string) *API {
	return &API{url: url}
}

func (a *API) Params(params H) *API {
	value := make(url.Values)
	for k, v := range params {
		value.Set(k, v)
	}
	a.url += "?" + value.Encode()
	return a
}

const codeforces = `https://codeforces.com/api`

func (a *API) Decode(x interface{}) error {
	response, err := http.Get(codeforces + a.url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(x)
}
