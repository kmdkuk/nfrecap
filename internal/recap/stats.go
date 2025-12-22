package recap

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/kmdkuk/nfrecap/internal/build"
)

type Stats struct {
	Year int

	TotalItems     int
	ItemsByMonth   [12]int
	ItemsByWeekday [7]int

	// streak
	LongestStreak      int
	LongestStreakStart *time.Time
	LongestStreakEnd   *time.Time
}

func ReadBuiltJSON(path string) (build.Built, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return build.Built{}, err
	}
	var x build.Built
	if err := json.Unmarshal(b, &x); err != nil {
		return build.Built{}, fmt.Errorf("invalid built json: %w", err)
	}
	return x, nil
}

func ComputeStats(built build.Built, year int) Stats {
	s := Stats{Year: year}

	// filter by year
	var dates []time.Time
	for _, it := range built.Items {
		d, err := time.Parse("2006-01-02", it.Date)
		if err != nil {
			continue
		}
		if d.Year() != year {
			continue
		}
		s.TotalItems++
		s.ItemsByMonth[int(d.Month())-1]++
		s.ItemsByWeekday[int(d.Weekday())]++
		dates = append(dates, d)
	}

	s.computeStreak(dates)
	return s
}

func (s *Stats) computeStreak(dates []time.Time) {
	if len(dates) == 0 {
		return
	}
	// unique day set
	seen := map[string]bool{}
	uniq := make([]time.Time, 0, len(dates))
	for _, d := range dates {
		k := d.Format("2006-01-02")
		if !seen[k] {
			seen[k] = true
			uniq = append(uniq, d)
		}
	}
	sort.Slice(uniq, func(i, j int) bool { return uniq[i].Before(uniq[j]) })

	bestLen := 1
	curLen := 1
	curStart := uniq[0]
	bestStart := uniq[0]
	bestEnd := uniq[0]

	for i := 1; i < len(uniq); i++ {
		prev := uniq[i-1]
		cur := uniq[i]
		if cur.Sub(prev) == 24*time.Hour {
			curLen++
		} else {
			curLen = 1
			curStart = cur
		}
		if curLen > bestLen {
			bestLen = curLen
			bestStart = curStart
			bestEnd = cur
		}
	}

	s.LongestStreak = bestLen
	s.LongestStreakStart = &bestStart
	s.LongestStreakEnd = &bestEnd
}
