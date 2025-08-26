package helper

import "math/rand"

func ToNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
