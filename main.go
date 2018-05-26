package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

var appExit chan bool

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

func init() {
	// Initalize exit channel
	appExit = make(chan bool)
}

func fatalf(str string, err error) {
	log.Criticalf(str, err)
	log.Flush()
	appExit <- true
}

func main() {
	// Read the configuration file
	conf, err := readConfig("config.json")

	if err != nil {
		fatalf("Error reading config file: %s", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go func() {
		r.Run()
	}()

	query := "Darker Shade"

	uri := fmt.Sprintf("%ssearch/index.xml?key=%s&q=%s", apiRoot, conf.GoodreadsAPIKey, url.QueryEscape(query))

	xml := getRequest(uri)

	response := &SearchResponse{}
	xmlUnmarshal(xml, response)

	fmt.Println(response)

	<-appExit
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		fatalf("Could not get the HTTP request: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fatalf("Could not get the HTTP Body: %s", err)
	}

	return body
}

func xmlUnmarshal(b []byte, i interface{}) {
	err := xml.Unmarshal(b, i)
	if err != nil {
		fatalf("Could not parse the returned XML: %s", err)
	}
}
