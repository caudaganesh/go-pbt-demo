package generator

import (
	"math/rand"
	"testing/quick"
	"time"
)

// GenerateQuickCheckConfig get maxcount, default is 15
func GenerateQuickCheckConfig(maxCount int) *quick.Config {
	if maxCount == 0 {
		maxCount = 15
	}
	return &quick.Config{
		MaxCount: maxCount,
		Rand:     rand.New(rand.NewSource(time.Now().Unix())),
	}
}
