package helper

import (
	"math/rand"
	"time"
)

const (
	min = 11111
	max = 99999
)

func DeliveryCode() int32 {
	// set seed
	rand.Seed(time.Now().UnixNano())
	// generate random number
	return int32(rand.Intn(max-min) + max)
}
