package utils

import (
	"math"
)

func ToFixed(x float64, p int) float64 {
	y := math.Pow(10, float64(p))
	return float64(int(x*y+math.Copysign(0.5, x*y))) / y
}
