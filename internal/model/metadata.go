package model

type Metadata struct {
	Provider string   `json:"provider"`
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year,omitempty"`
	Genres   []string `json:"genres,omitempty"`
	Runtime  int      `json:"runtime_min,omitempty"` // movie runtime or avg episode runtime
}
