package helper

func ToNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
