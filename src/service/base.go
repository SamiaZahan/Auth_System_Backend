package service

import "math/rand"

type Auth struct{}

func GenerateRandNum() int {
	min := 1000
	max := 9999
	return rand.Intn(max-min) + min
}
