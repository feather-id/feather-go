package feather

// ListMeta ...
// https://feather.id/docs/reference/api#pagination
type ListMeta struct {
	Objet      string `json:"Object"`
	URL        string `json:"url"`
	TotalCount uint32 `json:"total_count"`
}

// ListParams ...
// https://feather.id/docs/reference/api#pagination
type ListParams struct {
	Limit         *uint32 `json:"limit"`
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
}
