package utils

import (
	"math"
)

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func ToFixed(x float64, p int) float64 {
	y := math.Pow(10, float64(p))
	return float64(int(x*y+math.Copysign(0.5, x*y))) / y
}
