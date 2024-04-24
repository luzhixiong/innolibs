package utils

import (
	"math/rand"
	"time"
)

var rander *rand.Rand

func Rander() *rand.Rand {
	if rander == nil {
		rander = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return rander
}

// 浮点数随机数(min max 都要大于0)
func RandFloat32(min, max float32) float32 {
	if min >= max {
		return max
	}

	return min + rand.Float32()*(max-min)
}
