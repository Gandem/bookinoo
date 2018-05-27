package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"
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

func amazonSearchURL(query string) string {
	params := url.Values{
		"Operation":   []string{"ItemSearch"},
		"SearchIndex": []string{"Books"},
		"Keywords":    []string{query},
	}

	merged := mergeParams(params)

	url := url.URL{
		Scheme:   "https",
		Host:     appConfig.AmazonAPIRoot,
		Path:     "/onca/xml",
		RawQuery: merged.Encode(),
	}

	return url.String()
}

func mergeParams(extra url.Values) url.Values {
	params := url.Values{
		"AWSAccessKeyId": []string{appConfig.AmazonAccessKey},
		"AssociateTag":   []string{appConfig.AmazonAssociateID},
		"Service":        []string{"AWSECommerceService"},
		"Timestamp":      []string{time.Now().Format(time.RFC3339)},
	}
	for k, v := range extra {
		params[k] = v
	}

	// attach signature
	message := fmt.Sprintf("GET\n%s\n/onca/xml\n%s", appConfig.AmazonAPIRoot, strings.Replace(params.Encode(), "+", "%20", -1))
	signature := sign(message)
	params.Set("Signature", signature)

	return params
}

func sign(message string) string {
	key := []byte(appConfig.AmazonSecretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString((mac.Sum(nil)))
}
