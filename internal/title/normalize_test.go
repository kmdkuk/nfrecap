package title

import (
	"testing"

	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.NormalizedTitle
	}{
		{
			name:  "Movie basic",
			input: "Inception",
			expected: model.NormalizedTitle{
				RawTitle:  "Inception",
				WorkTitle: "Inception",
				Type:      "movie",
			},
		},
		{
			name:  "TV Show with Season",
			input: "Stranger Things: Season 1",
			expected: model.NormalizedTitle{
				RawTitle:  "Stranger Things: Season 1",
				WorkTitle: "Stranger Things",
				Type:      "tv",
				Season:    "Season 1",
			},
		},
		{
			name:  "TV Show with Season and Episode",
			input: "Black Mirror: Season 3: Nosedive",
			expected: model.NormalizedTitle{
				RawTitle:     "Black Mirror: Season 3: Nosedive",
				WorkTitle:    "Black Mirror",
				Type:         "tv",
				Season:       "Season 3",
				EpisodeTitle: "Nosedive",
			},
		},
		{
			name:  "Whitespace handling",
			input: "  Breaking Bad : Season 5  ",
			expected: model.NormalizedTitle{
				RawTitle:  "Breaking Bad : Season 5",
				WorkTitle: "Breaking Bad",
				Type:      "tv",
				Season:    "Season 5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
