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
