package helper

import (
	"math"
	"time"
)

func ProRateRemainingMonthsFraction(t time.Time) float64 {
	rem := 12 - int(t.Month()) // exclude current month
	return float64(rem) / 12.0
}

func ProRateRemainingMonthsPercent(t time.Time) float64 {
	return round2(ProRateRemainingMonthsFraction(t) * 100.0)
}

func round2(x float64) float64 { return math.Round(x*100) / 100 }
