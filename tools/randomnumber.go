package tools

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(limit int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit)
}
