package utils

import (
	"math/rand"
	"time"
)

func RandomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func RandomBool() bool {
	return RandomInt(2) > 0
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
