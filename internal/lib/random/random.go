package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	bytes := make([]rune, size)

	for i := range bytes {
		bytes[i] = chars[random.Intn(len(chars))]
	}

	return string(bytes)
}
