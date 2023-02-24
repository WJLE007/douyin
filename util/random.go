package util

import (
	"math/rand"
	"strings"
	"time"
)

func GenNumCode(length int) string {
	source := "0123456789"
	rand.Seed(time.Now().UnixNano())
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(source[rand.Intn(len(source))])
	}
	return builder.String()
}

func GenStrCode(length int) string {
	source := "0123456789abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().UnixNano())
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteByte(source[rand.Intn(len(source))])
	}
	return builder.String()
}
