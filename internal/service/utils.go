package service

import (
	"math"
)

func toCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func fromCents(cents int64) float64 {
	return float64(cents) / 100.0
}
