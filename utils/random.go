package utils

import (
	"math/rand"
	"strings"
	"time"
)

var alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// rand.Seed(time.Now().Unix())
	rand.NewSource(time.Now().Unix())
}

func generateRandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func generateRandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func GenerateRandomOwner() string {
	return generateRandomString(6)
}

func GenerateRandomCurrency() string {
	currencies := []string{"INR", "USD", "CAD"}
	k := len(currencies)
	return currencies[rand.Intn(k)]
}

func GenerateRandomMoney() int64 {
	return generateRandomInt(1, 1000)
}
