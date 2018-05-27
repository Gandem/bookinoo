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

var (
	appExit   chan bool
	appConfig *apiConfig
)

type AmazonItem struct {
	ASIN             string
	ParentASIN       string
	DetailPageURL    string
	SalesRank        string
	ItemLinks        []AmazonItemLink `xml:"ItemLinks>ItemLink"`
	SmallImage       AmazonImage
	MediumImage      AmazonImage
	LargeImage       AmazonImage
	ImageSets        []AmazonImageSet `xml:"ImageSets>ImageSet"`
	ItemAttributes   AmazonItemAttributes
	EditorialReviews []AmazonEditorialReview `xml:"EditorialReviews>EditorialReview"`
}

type AmazonEditorialReview struct {
	Source  string
	Content string
}

type AmazonItemLink struct {
	Description string
	URL         string
}

type AmazonImageSet struct {
	Category       string `xml:"Category,attr"`
	SwatchImage    AmazonImage
	SmallImage     AmazonImage
	ThumbnailImage AmazonImage
	TinyImage      AmazonImage
	MediumImage    AmazonImage
	LargeImage     AmazonImage
}

type AmazonImage struct {
	URL    string
	Height uint16
	Width  uint16
}

type AmazonItemAttributes struct {
	Title     string
	Brand     string
	ListPrice AmazonItemPrice
}

type AmazonItemPrice struct {
	Amount         int64
	CurrencyCode   string
	FormattedPrice string
}

type AmazonItems struct {
	Items []AmazonItem `xml:"Item"`
}

type AmazonItemSearchResponse struct {
	XMLName     xml.Name    `xml:"ItemSearchResponse"`
	AmazonItems AmazonItems `xml:"Items"`
}

// GoodreadsAuthor represents an author as returned by the Goodreads API
type GoodreadsAuthor struct {
	GoodreadsID int    `xml:"id"`
	Name        string `xml:"name"`
}

// GoodreadsBook represents a book as returned by the Goodreads API
type GoodreadsBook struct {
	GoodreadsID int             `xml:"id"`
	Title       string          `xml:"title"`
	Author      GoodreadsAuthor `xml:"author"`
	ImageURL    string          `xml:"image_url"`
}

// SearchResponse represents a response to a search query as returned by the Goodreads API
type GoodreadsSearchResponse struct {
	Books []GoodreadsBook `xml:"search>results>work>best_book"`
}

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

		uri := fmt.Sprintf("https://%s/search/index.xml?key=%s&q=%s", appConfig.GoodreadsAPIRoot, appConfig.GoodreadsAPIKey, url.QueryEscape(query))

		goodreadsXML := getRequest(uri)
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
