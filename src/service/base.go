package service

import "math/rand"

func GenerateRandNum() int {
	min := 1000
	max := 9999
	return rand.Intn(max-min) + min
}
