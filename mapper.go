package feather

// String returns a pointer to the provided string value.
func String(v string) *string {
	return &v
}

// Bool returns a pointer to the provided bool value.
func Bool(v bool) *bool {
	return &v
}

// UInt32 returns a pointer to the provided uint32 value.
func UInt32(v uint32) *uint32 {
	return &v
}
