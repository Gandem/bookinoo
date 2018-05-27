package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

var (
	appExit   chan bool
	appConfig *apiConfig
)

func init() {
	// Initalize exit channel
	appExit = make(chan bool)

	var err error
	// Read the configuration file
	appConfig, err = readConfig("config.json")
	if err != nil {
		fatalf("Error reading config file: %s", err)
	}
}

func fatalf(str string, err error) {
	log.Criticalf(str, err)
	log.Flush()
	appExit <- true
}

func main() {
	r := gin.Default()
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")

		goodreadsXML := getRequest(goodreadsSearchURL(query))
		amazonXML := getRequest(amazonSearchURL(query))

		goodreadsResponse := &GoodreadsSearchResponse{}
		xmlUnmarshal(goodreadsXML, goodreadsResponse)

		amazonResponse := &AmazonItemSearchResponse{}
		xmlUnmarshal(amazonXML, amazonResponse)

		c.JSON(200, gin.H{
			"goodreads": goodreadsResponse,
			"amazon":    amazonResponse,
		})
	})
	go func() {
		r.Run()
	}()

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
