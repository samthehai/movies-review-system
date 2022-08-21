package utils

import "time"

func StringPtr(v string) *string {
	return &v
}

func TimePtr(v time.Time) *time.Time {
	return &v
}

func Int64Ptr(v int64) *int64 {
	return &v
}

func Uint64Ptr(v uint64) *uint64 {
	return &v
}
