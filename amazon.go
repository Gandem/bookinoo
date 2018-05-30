package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type AmazonSearchItem struct {
	ISBN   string      `xml:"ASIN"`
	Title  string      `xml:"ItemAttributes>Title"`
	Author string      `xml:"ItemAttributes>Author"`
	Image  AmazonImage `xml:"MediumImage"`
}

type AmazonImage struct {
	URL    string
	Height uint16
	Width  uint16
}

type AmazonSearchResponse struct {
	Items []AmazonSearchItem `xml:"Items>Item"`
}

func amazonSearchURL(query string) string {
	params := url.Values{
		"Operation":     []string{"ItemSearch"},
		"SearchIndex":   []string{"Books"},
		"Keywords":      []string{query},
		"ResponseGroup": []string{"Images,ItemAttributes"},
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

type AmazonReviewItem struct {
	URL     string `xml:"Items>Item>DetailPageURL"`
	Reviews string `xml:"Items>Item>CustomerReviews>IFrameURL"`
}

type AmazonReviewBackend struct {
	name string
}

func (a AmazonReviewBackend) getName() string {
	return a.name
}

func (AmazonReviewBackend) query(isbn string) reviewsGroup {
	params := url.Values{
		"Operation":     []string{"ItemLookup"},
		"SearchIndex":   []string{"Books"},
		"ItemId":        []string{isbn},
		"IdType":        []string{"ISBN"},
		"ResponseGroup": []string{"Large"},
	}

	merged := mergeParams(params)

	url := url.URL{
		Scheme:   "https",
		Host:     appConfig.AmazonAPIRoot,
		Path:     "/onca/xml",
		RawQuery: merged.Encode(),
	}

	amazonResponse := &AmazonReviewItem{}
	// TODO:
	amazonReviewsXML, _ := getRequest(url.String())
	xmlUnmarshal(amazonReviewsXML, amazonResponse)

	return reviewsGroup{
		Platform:      "amazon",
		URL:           amazonResponse.URL,
		Reviews:       amazonResponse.Reviews,
		ReviewsType:   "url",
		AverageRating: "-1",
	}
}
