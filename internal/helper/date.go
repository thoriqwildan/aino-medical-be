package helper

import "time"

func isAnniversary(join, now time.Time) bool {
	if join.IsZero() {
		return false
	}
	if now.Location() != join.Location() {
		now = now.In(join.Location())
	}
	if now.Before(join) {
		return false
	}

	jY, jM, jD := join.Date()
	nY, nM, nD := now.Date()

	if jM == time.February && jD == 29 {
		isLeap := func(y int) bool {
			return (y%4 == 0 && y%100 != 0) || (y%400 == 0)
		}
		if !isLeap(nY) {
			return nM == time.February && nD == 28 && (jY != nY && jY < nY)
		}
	}

	return nM == jM && nD == jD && (jY != nY && jY < nY)
}
