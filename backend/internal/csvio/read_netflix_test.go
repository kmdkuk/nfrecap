package csvio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNetflixCSV(t *testing.T) {
	// Refined structure
	t.Run("Normal Case", func(t *testing.T) {
		input := `Title,Date
"Stranger Things: Season 1: Chapter One",1/1/23
"Inception",12/31/22`

		recs, err := ParseNetflixCSV(strings.NewReader(input))
		require.NoError(t, err)
		require.Len(t, recs, 2)

		assert.Equal(t, "Stranger Things: Season 1: Chapter One", recs[0].Title)
		assert.Equal(t, "2023-01-01", recs[0].Date.Format("2006-01-02"))

		assert.Equal(t, "Inception", recs[1].Title)
		assert.Equal(t, "2022-12-31", recs[1].Date.Format("2006-01-02"))
	})

	t.Run("Empty CSV", func(t *testing.T) {
		input := "" // ReadAll returns empty rows
		_, err := ParseNetflixCSV(strings.NewReader(input))
		assert.ErrorContains(t, err, "empty csv")
	})

	t.Run("Header Only", func(t *testing.T) {
		input := "Title,Date"
		recs, err := ParseNetflixCSV(strings.NewReader(input))
		assert.NoError(t, err)
		assert.Len(t, recs, 0)
	})

	t.Run("Short Row", func(t *testing.T) {
		input := `Title,Date
Inception`
		recs, err := ParseNetflixCSV(strings.NewReader(input))
		assert.NoError(t, err)
		assert.Len(t, recs, 0) // Should skip short row
	})

	t.Run("Invalid Date", func(t *testing.T) {
		input := `Title,Date
Inception,2023-01-01` // Wrong format, expects M/D/YY (1/1/23)
		_, err := ParseNetflixCSV(strings.NewReader(input))
		assert.ErrorContains(t, err, "date parse failed")
	})

	t.Run("Date with different single digit format", func(t *testing.T) {
		// 1/1/23 vs 01/01/23
		input := `Title,Date
Movie A,1/1/23
Movie B,01/01/23`
		recs, err := ParseNetflixCSV(strings.NewReader(input))
		require.NoError(t, err)
		require.Len(t, recs, 2)
		assert.Equal(t, "2023-01-01", recs[0].Date.Format("2006-01-02"))
		assert.Equal(t, "2023-01-01", recs[1].Date.Format("2006-01-02"))
	})
}
