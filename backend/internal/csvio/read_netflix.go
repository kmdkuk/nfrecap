package csvio

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kmdkuk/nfrecap/internal/model"
)

func ReadNetflixCSV(path string) ([]model.ViewingRecord, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseNetflixCSV(f)
}

func ParseNetflixCSV(r io.Reader) ([]model.ViewingRecord, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1

	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty csv")
	}

	out := make([]model.ViewingRecord, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue // header
		}
		if len(row) < 2 {
			continue
		}
		title := strings.TrimSpace(row[0])
		ds := strings.TrimSpace(row[1])

		// Netflix viewing history often uses M/D/YY like "12/13/25"
		d, err := time.Parse("1/2/06", ds)
		if err != nil {
			return nil, fmt.Errorf("date parse failed at line %d: %q: %w", i+1, ds, err)
		}

		out = append(out, model.ViewingRecord{
			Title: title,
			Date:  time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local),
		})
	}
	return out, nil
}
