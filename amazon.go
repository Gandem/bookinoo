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
