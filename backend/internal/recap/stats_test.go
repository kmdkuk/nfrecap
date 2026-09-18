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
		{
			name: "Genre Ranking Top/Worst",
			year: 2023,
			items: []build.BuiltItem{
				{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "A", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 8.0}},
				{Date: "2023-01-02", Normalized: model.NormalizedTitle{WorkTitle: "B", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 2.0}},
				{Date: "2023-01-03", Normalized: model.NormalizedTitle{WorkTitle: "C", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 9.0}},
				{Date: "2023-01-04", Normalized: model.NormalizedTitle{WorkTitle: "D", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 3.0}},
				{Date: "2023-01-05", Normalized: model.NormalizedTitle{WorkTitle: "E", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 5.0}},
				// Unrated one, should be ignored
				{Date: "2023-01-06", Normalized: model.NormalizedTitle{WorkTitle: "X", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Action"}, VoteAverage: 0.0}},
			},
			expected: func(t *testing.T, s Stats) {
				// Sorted by Desc: C(9), A(8), E(5), D(3), B(2)
				// Top 3: C, A, E
				if assert.Contains(t, s.GenreTopWorks, "Action") {
					top := s.GenreTopWorks["Action"]
					if assert.Len(t, top, 3) {
						assert.Equal(t, "C", top[0].Title)
						assert.Equal(t, "A", top[1].Title)
						assert.Equal(t, "E", top[2].Title)
					}
				}

				// Worst 3: B(2), D(3), E(5) (Lowest first)
				if assert.Contains(t, s.GenreWorstWorks, "Action") {
					worst := s.GenreWorstWorks["Action"]
					if assert.Len(t, worst, 3) {
						assert.Equal(t, "B", worst[0].Title)
						assert.Equal(t, "D", worst[1].Title)
						assert.Equal(t, "E", worst[2].Title)
					}
				}
			},
		},
		{
			name: "Genre Sample Movies Sorted by Popularity",
			year: 2023,
			items: []build.BuiltItem{
				{Date: "2023-01-01", Normalized: model.NormalizedTitle{WorkTitle: "Low Pop", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 10.0}},
				{Date: "2023-01-02", Normalized: model.NormalizedTitle{WorkTitle: "High Pop", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 100.0}},
				{Date: "2023-01-03", Normalized: model.NormalizedTitle{WorkTitle: "Mid Pop", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 50.0}},
				{Date: "2023-01-04", Normalized: model.NormalizedTitle{WorkTitle: "Super Pop", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 200.0}},
				{Date: "2023-01-05", Normalized: model.NormalizedTitle{WorkTitle: "Tiny Pop", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 1.0}},
				{Date: "2023-01-06", Normalized: model.NormalizedTitle{WorkTitle: "Extra", Type: "movie"}, Metadata: &model.Metadata{Genres: []string{"Drama"}, Popularity: 5.0}},
			},
			expected: func(t *testing.T, s Stats) {
				if assert.Contains(t, s.GenreSampleMovies, "Drama") {
					samples := s.GenreSampleMovies["Drama"]
					// Should be top 5 by popularity desc
					// Items: Super(200), High(100), Mid(50), Low(10), Extra(5)|Tiny(1)
					// Expected Top 5: Super, High, Mid, Low, Extra (Tiny is 6th, dropped)
					if assert.Len(t, samples, 5) {
						assert.Equal(t, "Super Pop", samples[0].Title)
						assert.Equal(t, "High Pop", samples[1].Title)
						assert.Equal(t, "Mid Pop", samples[2].Title)
						assert.Equal(t, "Low Pop", samples[3].Title)
						assert.Equal(t, "Extra", samples[4].Title)
					}
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
