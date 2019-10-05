package generator

import (
	"math/rand"
	"time"
)

func initRand() *rand.Rand {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r
}
