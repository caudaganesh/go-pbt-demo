package generator

import (
	"math/rand"
	"time"
)

// RandomTimestamp return random time
func GenerateRandomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
	randomNow := time.Unix(randomTime, 0)

	return randomNow
}
