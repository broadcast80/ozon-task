package domain

import "time"

// ???
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type Link struct {
	ID        string
	Alias     string
	URL       string
	CreatedAt time.Time
}
