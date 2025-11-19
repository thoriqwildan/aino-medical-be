package tests

import (
	"testing"
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/helper"
)

func TestCalculateEmployeeProrate(t *testing.T) {
	prorateNow := helper.ProRateRemainingMonthsPercent(time.Now(), time.Now().AddDate(0, -2, 0))
	join, err := time.Parse("2006-01-02", "2025-08-29")
	if err != nil {
		t.Error(err)
	}
	joinDate := helper.CustomDate(join)
	prorateJoinDate := helper.ProRateRemainingMonthsPercent(time.Now(), time.Time(joinDate))

	t.Log(prorateNow)
	t.Log(prorateJoinDate)
}
