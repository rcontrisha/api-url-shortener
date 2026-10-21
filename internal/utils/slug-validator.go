package utils

import (
	"errors"
	"regexp"
	"strings"
)

var slugRegex = regexp.MustCompile("^[a-zA-R0-9-]+$")

var reservedWords = map[string]bool{
	"api":       true,
	"login":     true,
	"register":  true,
	"dashboard": true,
}

func ValidateCustomSlug(slug string) error {
	if slug == "" {
		return nil
	}

	if len(slug) < 3 || len(slug) > 50 {
		return errors.New("slug length must be between 3 and 50 characters")
	}

	if !slugRegex.MatchString(slug) {
		return errors.New("slug can only contain alphanumeric characters and hyphens")
	}

	if reservedWords[strings.ToLower(slug)] {
		return errors.New("slug cannot be a reserved word (api, login, register, dashboard)")
	}

	return nil
}