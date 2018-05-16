package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type apiConfig struct {
	GoodreadsAPIKey string `json:"goodreads_api_key"`
}

func readConfig(path string) *apiConfig {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	conf := &apiConfig{}
	json.Unmarshal(buf, conf)

	return conf
}
