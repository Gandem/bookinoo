package main

import (
	"encoding/json"
	"io/ioutil"
)

type apiConfig struct {
	GoodreadsAPIKey   string `json:"goodreads_api_key"`
	GoodreadsAPIRoot  string `json:"goodreads_api_rooturl"`
	AmazonAccessKey   string `json:"amazon_access_key"`
	AmazonSecretKey   string `json:"amazon_secret_key"`
	AmazonAssociateID string `json:"amazon_associate_id"`
	AmazonAPIRoot     string `json:"amazon_api_rooturl"`
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
