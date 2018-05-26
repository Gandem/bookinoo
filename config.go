package main

import (
	"encoding/json"
	"io/ioutil"
)

type apiConfig struct {
	GoodreadsAPIKey string `json:"goodreads_api_key"`
}

func readConfig(path string) (*apiConfig, error) {
	conf := &apiConfig{}

	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return conf, err
	}

	json.Unmarshal(buf, conf)

	return conf, nil
}
