package feather

// ListMeta ...
// https://feather.id/docs/reference/api#pagination
type ListMeta struct {
	HasMore    bool   `json:"has_more"`
	TotalCount uint32 `json:"total_count"`
	URL        string `json:"url"`
}

// ListParams ...
// https://feather.id/docs/reference/api#pagination
type ListParams struct {
	Limit         *uint32 `json:"limit"`
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
}
