package helper

import (
	"fmt"
	"time"
)

// CustomDate hanya untuk mem-parsing string "YYYY-MM-DD"
type CustomDate time.Time

// MarshalJSON mengubah CustomDate menjadi string "YYYY-MM-DD" untuk respons
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02"))), nil
}

// UnmarshalJSON mengubah string "YYYY-MM-DD" dari request menjadi CustomDate
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	s := string(data)
	// Hapus tanda kutip jika ada
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	if s == "null" {
		*cd = CustomDate(time.Time{})
		return nil
	}

	// Parsing dengan format "2006-01-02"
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format for YYYY-MM-DD: %w", err)
	}
	*cd = CustomDate(t)
	return nil
}