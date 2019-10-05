package generator

import (
	"math/rand"
)

func GenerateRandomInts() []int {
	maxLen := initRand().Intn(20)
	res := []int{}
	for index := 0; index < maxLen; index++ {
		res = append(res, rand.Int())
	}
	return res
}

func GenerateRandomInt() int {
	return rand.Int()
}
