package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	apiRoot = "http://www.goodreads.com/"
)

// GoodreadsAuthor represents an author as returned by the Goodreads API
type GoodreadsAuthor struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

// GoodreadsBook represents a book as returned by the Goodreads API
type GoodreadsBook struct {
	ID       int             `xml:"id"`
	Title    string          `xml:"title"`
	Author   GoodreadsAuthor `xml:"author"`
	ImageURL string          `xml:"image_url"`
}

// SearchResponse represents a response to a search query as returned by the Goodreads API
type SearchResponse struct {
	Books []GoodreadsBook `xml:"search>results>work>best_book"`
}

func main() {
	conf := readConfig("config.json")

	query := "Darker Shade"

	uri := fmt.Sprintf("%ssearch/index.xml?key=%s&q=%s", apiRoot, conf.GoodreadsAPIKey, url.QueryEscape(query))

	xml := getRequest(uri)

	response := &SearchResponse{}
	xmlUnmarshal(xml, response)

	fmt.Println(response)
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func xmlUnmarshal(b []byte, i interface{}) {
	err := xml.Unmarshal(b, i)
	if err != nil {
		log.Fatal(err)
	}
}
