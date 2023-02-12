package util

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomNumber() string {
	return strconv.Itoa(rand.New(
		rand.NewSource(time.Now().UnixNano())).Int())
}
