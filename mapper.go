package feather

import "time"

// String returns a pointer to the provided string value.
func String(v string) *string {
	return &v
}

// UInt32 returns a pointer to the provided uint32 value.
func UInt32(v uint32) *uint32 {
	return &v
}

// Time returns a pointer to the provided time.Time value.
func Time(v time.Time) *time.Time {
	return &v
}
