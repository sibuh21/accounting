package service

func toCents(amount float64) int64 {
	return int64(amount*100 + 0.5)
}

func fromCents(cents int64) float64 {
	return float64(cents) / 100.0
}
