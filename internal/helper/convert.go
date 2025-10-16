package helper

import "time"

func PtrFloat64(v float64) *float64 { return &v }
func PtrString(s string) *string    { return &s }

func PtrDate(date time.Time) *time.Time { return &date }