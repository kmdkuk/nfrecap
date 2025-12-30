package recap

import (
	"testing"

	"github.com/kmdkuk/nfrecap/internal/build"
	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestComputeStats(t *testing.T) {
	// Helper to create date string
	date := func(s string) string { return s }

	tests := []struct {
		name     string
		year     int
		items    []build.BuiltItem
		expected func(*testing.T, Stats)
	}{
		{
			name: "Basic Stats Calculation",
			year: 2023,
			items: []build.BuiltItem{
				{
					Date: date("2023-01-01"),
					Normalized: model.NormalizedTitle{
						WorkTitle: "Movie A",
						Type:      "movie",
					},
					Metadata: &model.Metadata{
						Runtime: 120,
						Genres:  []string{"Action"},
					},
				},
				{
					Date: date("2023-01-02"),
					Normalized: model.NormalizedTitle{
						WorkTitle: "Movie B",
						Type:      "movie",
					},
					Metadata: &model.Metadata{
						Runtime: 90,
						Genres:  []string{"Action", "Comedy"},
					},
				},
			},
			expected: func(t *testing.T, s Stats) {
				assert.Equal(t, 2, s.TotalViews)
				assert.Equal(t, 210, s.TotalDurationMin)
				assert.Equal(t, 2, s.ActiveDays)
				assert.Equal(t, 2023, s.Year)

				// Genre Stats
				// Action: 2 views, 120+90=210 min
				// Comedy: 1 view, 90 min
				var action, comedy GenreStat
				for _, g := range s.GenreStats {
					if g.Name == "Action" {
						action = g
					}
					if g.Name == "Comedy" {
						comedy = g
					}
				}
				assert.Equal(t, "Action", action.Name)
				assert.Equal(t, 210, action.DurationMin)
				assert.Equal(t, "Comedy", comedy.Name)
				assert.Equal(t, 90, comedy.DurationMin)
			},
		},
		{
			name: "Streak Calculation",
			year: 2023,
			items: []build.BuiltItem{
				{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "A", Type: "movie"}, Metadata: &model.Metadata{Runtime: 60}},
				{Date: "2023-01-02", Normalized: model.NormalizedTitle{WorkTitle: "B", Type: "movie"}, Metadata: &model.Metadata{Runtime: 60}},
				{Date: "2023-01-03", Normalized: model.NormalizedTitle{WorkTitle: "C", Type: "movie"}, Metadata: &model.Metadata{Runtime: 60}},
				// Gap 01-04 (missing 04) -> Wait, 03 to 05 is gap of 1 day (04).
				{Date: "2023-01-05", Normalized: model.NormalizedTitle{WorkTitle: "D", Type: "movie"}, Metadata: &model.Metadata{Runtime: 60}},
			},
			expected: func(t *testing.T, s Stats) {
				assert.Equal(t, 4, s.TotalViews)
				// Streak: 1,2,3 -> 3 days
				// Streak: 5 -> 1 day
				// Top streak should be 3
				if assert.NotEmpty(t, s.TopStreaks) {
					assert.Equal(t, 3, s.TopStreaks[0].Days)
					assert.Equal(t, "2023-01-01", s.TopStreaks[0].Start.Format("2006-01-02"))
					assert.Equal(t, "2023-01-03", s.TopStreaks[0].End.Format("2006-01-02"))
				}

				// Max Gap
				// 01-03 contiguous.
				// 03 to 05. Difference is 2 days. Gap is 1 day (04).
				assert.Equal(t, 1, s.MaxGap.Days)
				assert.Equal(t, "2023-01-04", s.MaxGap.Start.Format("2006-01-02"))
			},
		},
		{
			name: "Filtering by Year",
			year: 2023,
			items: []build.BuiltItem{
				{Date: "2022-12-31", Normalized: model.NormalizedTitle{WorkTitle: "Old"}, Metadata: &model.Metadata{Runtime: 10}},
				{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "New"}, Metadata: &model.Metadata{Runtime: 10}},
			},
			expected: func(t *testing.T, s Stats) {
				assert.Equal(t, 1, s.TotalViews)
				if assert.Len(t, s.TopTitlesByDuration, 1) {
					assert.Equal(t, "New", s.TopTitlesByDuration[0].Title)
				}
			},
		},
		{
			name: "Unresolved Items",
			year: 2023,
			items: []build.BuiltItem{
				{
					Date: "2023-01-01",
					Normalized: model.NormalizedTitle{
						WorkTitle: "Unknown Title",
						Type:      "movie",
					},
					Metadata: nil, // Unresolved
				},
			},
			expected: func(t *testing.T, s Stats) {
				assert.Equal(t, 1, s.TotalViews) // It counts as a view even if unresolved? Implementation check: yes
				assert.Equal(t, 1, s.UnresolvedCount)
				if assert.Len(t, s.UnresolvedList, 1) {
					assert.Equal(t, "Unknown Title", s.UnresolvedList[0].Title)
					assert.Equal(t, 1, s.UnresolvedList[0].Views)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := build.Built{
				Items: tt.items,
			}
			got := ComputeStats(input, tt.year)
			tt.expected(t, got)
		})
	}
}

func TestComputeStreaksAndGaps_EdgeCases(t *testing.T) {
	// Focusing specifically on logic inside computeStreaksAndGaps via ComputeStats
	// using empty or single item

	t.Run("Empty", func(t *testing.T) {
		s := ComputeStats(build.Built{}, 2023)
		assert.Empty(t, s.TopStreaks)
		assert.Equal(t, 0, s.MaxGap.Days)
	})

	t.Run("Single Day", func(t *testing.T) {
		items := []build.BuiltItem{
			{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "A"}, Metadata: &model.Metadata{}},
			{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "B"}, Metadata: &model.Metadata{}},
		}
		s := ComputeStats(build.Built{Items: items}, 2023)
		assert.Equal(t, 2, s.TotalViews)
		assert.Equal(t, 1, s.ActiveDays)

		if assert.Len(t, s.TopStreaks, 1) {
			assert.Equal(t, 1, s.TopStreaks[0].Days)
		}
		assert.Equal(t, 0, s.MaxGap.Days)
	})
}
