package main

import (
	"fmt"
	"net/url"
)

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

// GoodreadsSearchResponse represents a response to a search query as returned by the Goodreads API
type GoodreadsSearchResponse struct {
	Books []GoodreadsBook `xml:"search>results>work>best_book"`
}

func goodreadsSearchURL(query string) string {
	return fmt.Sprintf("https://%s/search/index.xml?key=%s&q=%s",
		appConfig.GoodreadsAPIRoot,
		appConfig.GoodreadsAPIKey,
		url.QueryEscape(query))
}

type GoodreadsReviewItem struct {
	URL           string `xml:"book>url"`
	Reviews       string `xml:"book>reviews_widget"`
	AverageRating string `xml:"book>average_rating"`
}
type GoodreadsReviewBackend struct{}

func (GoodreadsReviewBackend) query(isbn string) reviewsGroup {
	uri := fmt.Sprintf("https://%s/book/isbn/%s?key=%s",
		appConfig.GoodreadsAPIRoot,
		url.QueryEscape(isbn),
		appConfig.GoodreadsAPIKey,
	)

	goodreadsResponse := &GoodreadsReviewItem{}
	goodreadsReviewXML, statusCode := getRequest(uri)

	if statusCode == 404 {
		return reviewsGroup{}
	}

	xmlUnmarshal(goodreadsReviewXML, goodreadsResponse)

	return reviewsGroup{
		Platform:      "goodreads",
		URL:           goodreadsResponse.URL,
		Reviews:       goodreadsResponse.Reviews,
		ReviewsType:   "html",
		AverageRating: goodreadsResponse.AverageRating,
	}
}
