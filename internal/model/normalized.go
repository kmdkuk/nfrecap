package model

type NormalizedTitle struct {
	RawTitle     string `json:"raw_title"`
	WorkTitle    string `json:"work_title"`
	Type         string `json:"type"` // "movie" | "tv" | "unknown"
	Season       string `json:"season,omitempty"`
	EpisodeTitle string `json:"episode_title,omitempty"`
}
