package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/contrib/gin-gonic/gintrace"
	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

type reviewsGroup struct {
	Platform      string
	URL           string
	Reviews       string
	ReviewsType   string
	AverageRating string
}

type reviewBackend interface {
	query(isbn string) reviewsGroup
}

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

	r.Use(gintrace.Middleware("bookinoo"))

	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")

		// TODO:
		amazonQuerySpan := tracer.NewChildSpanFromContext("amazon.query", c.Request.Context())
		amazonXML, _ := getRequest(amazonSearchURL(query))
		amazonQuerySpan.Finish()

		amazonResponse := &AmazonSearchResponse{}

		xmlUnmarshalSpan := tracer.NewChildSpanFromContext("xml.unmarshal", c.Request.Context())
		xmlUnmarshal(amazonXML, amazonResponse)
		xmlUnmarshalSpan.Finish()

		c.JSON(200, gin.H{
			"books": amazonResponse.Items,
		})
	})

	reviewBackends := [2]reviewBackend{
		AmazonReviewBackend{},
		GoodreadsReviewBackend{},
	}

	r.GET("/reviews", func(c *gin.Context) {
		isbn := c.Query("ISBN")
		reviewsChan := make(chan reviewsGroup)
		wg := sync.WaitGroup{}

		for _, backend := range reviewBackends {
			wg.Add(1)
			go func(backend reviewBackend) {
				reviewsChan <- backend.query(isbn)
				wg.Done()
			}(backend)
		}
		reviews := make([]reviewsGroup, 0)

		done := make(chan bool)
		go func() {
			for rg := range reviewsChan {
				if rg != (reviewsGroup{}) {
					reviews = append(reviews, rg)
				}
			}
			done <- true
		}()
		wg.Wait()
		close(reviewsChan)
		<-done

		c.JSON(200, gin.H{
			"reviews": reviews,
		})
	})

	go func() {
		r.Run()
	}()

	<-appExit
}

func getRequest(uri string) ([]byte, int) {
	res, err := http.Get(uri)
	if err != nil {
		fatalf("Could not get the HTTP request: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fatalf("Could not get the HTTP Body: %s", err)
	}

	return body, res.StatusCode
}

func xmlUnmarshal(b []byte, i interface{}) {
	err := xml.Unmarshal(b, i)
	if err != nil {
		fatalf("Could not parse the returned XML: %s", err)
	}
}
