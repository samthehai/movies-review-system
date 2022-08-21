package utils

import "time"

func MustRFC3339Time(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		panic(err)
	}

	return t
}
