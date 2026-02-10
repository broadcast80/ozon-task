package models

import "errors"

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

var ErrDuplicate = errors.New("duplicate url")
var ErrNotFound = errors.New("no such url")
