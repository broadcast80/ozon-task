package utils

import (
	"math/rand"
	"time"
)

func Encode(size int, charset string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}

	return string(b)
}

// может быть проблема с дублями
