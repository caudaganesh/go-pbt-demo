package generator

import (
	"errors"
	"math/rand"
)

var randError = map[int]error{
	0: errors.New("some error"),
	1: nil,
}

func GenerateRandomError() error {
	return randError[rand.Intn(2)]
}
