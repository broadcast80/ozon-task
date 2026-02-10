package models

import "errors"

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

type ResponseURL struct {
	URL string `json:"url"`
}

type ResponseAlias struct {
	Alias string `json:"alias"`
}

var ErrDuplicate = errors.New("duplicate url")
