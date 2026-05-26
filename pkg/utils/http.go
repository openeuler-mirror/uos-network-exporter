package utils

import (
	"fmt"
	"net/url"
)

func ValidateURI(uri string) error {
	parsedURL, err := url.ParseRequestURI(uri)
