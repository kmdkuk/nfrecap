package model

type Metadata struct {
	Provider    string   `json:"provider"`
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Year        int      `json:"year,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Runtime     int      `json:"runtime_min,omitempty"` // movie runtime or avg episode runtime
	VoteAverage float64  `json:"vote_average,omitempty"`
	Popularity  float64  `json:"popularity,omitempty"`
	PosterPath  string   `json:"poster_path,omitempty"`
}
