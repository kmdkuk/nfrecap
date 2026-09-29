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

	// Metadata from Built
	GeneratedAt string
	SourceFile  string

	// Basic Stats
	TotalViews       int
	TotalDurationMin int
	ActiveDays       int

	// Streaks & Gaps
	TopStreaks []Streak
	MaxGap     Gap

	// Monthly / Weekday
	MonthlyStats map[time.Month]Metric
	WeekdayStats map[time.Weekday]Metric

	// Genres
	GenreStats        []GenreStat
	GenreMonthSpike   map[string]Spike       // Genre -> Spike Info
	GenreSampleMovies map[string][]TitleStat // Genre -> List of movie sample TitleStats
	GenreTopWorks     map[string][]TitleStat // Genre -> Top 3 Works (by vote average)
	GenreWorstWorks   map[string][]TitleStat // Genre -> Worst 3 Works (by vote average)

	// Titles
	TopTitlesByDuration []TitleStat
	TopTitlesByViews    []TitleStat

	// TV Series
	TopSeriesByDuration []SeriesStat
	TopSeriesByViews    []SeriesStat

	// Unresolved
	UnresolvedCount int
	UnresolvedList  []UnresolvedItem
}

type Metric struct {
	Views       int
	DurationMin int
}

type Streak struct {
	Days  int
	Start time.Time
	End   time.Time
}

type Gap struct {
	Days  int
	Start time.Time
	End   time.Time
}

type GenreStat struct {
	Name        string
	DurationMin int
	Views       int
	Share       float64
}

type Spike struct {
	Month       time.Month
	DurationMin int
}

type TitleStat struct {
	Title       string
	Type        string // movie or tv
	DurationMin int
	Views       int
	VoteAverage float64
	PosterPath  string
	Popularity  float64
}

type SeriesStat struct {
	SeriesName  string
	DurationMin int
	Views       int
	SpanStart   time.Time
	SpanEnd     time.Time
}

type UnresolvedItem struct {
	Title string
	Type  string
	Views int
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
	s := Stats{
		Year:         year,
		GeneratedAt:  built.GeneratedAt,
		SourceFile:   built.Source,
		MonthlyStats: make(map[time.Month]Metric),
		WeekdayStats: make(map[time.Weekday]Metric),

		GenreSampleMovies: make(map[string][]TitleStat),
		GenreTopWorks:     make(map[string][]TitleStat),
		GenreWorstWorks:   make(map[string][]TitleStat),
	}

	// Internal aggregation maps
	genreMap := make(map[string]*Metric)                 // Genre -> Metric
	genreMonthMap := make(map[string]map[time.Month]int) // Genre -> Month -> Duration
	titleMap := make(map[string]*TitleStat)              // "Title|Type" -> TitleStat
	titleGenres := make(map[string][]string)             // "Title|Type" -> Genres
	seriesMap := make(map[string]*SeriesStat)            // SeriesName -> SeriesStat
	unresolvedMap := make(map[string]int)                // Title|Type -> count
	genreSeen := make(map[string]map[string]bool)        // Genre -> Title -> bool (for dedupe)

	var dates []time.Time

	for _, it := range built.Items {
		d, err := time.Parse("2006-01-02", it.Date)
		if err != nil {
			continue
		}
		if d.Year() != year {
			continue
		}

		dates = append(dates, d)

		s.TotalViews++

		// Metadata handling
		dur := 0
		var genres []string

		if it.Metadata != nil {
			dur = it.Metadata.Runtime
			genres = it.Metadata.Genres
		} else {
			// Unresolved
			key := fmt.Sprintf("%s|%s", it.Normalized.WorkTitle, it.Normalized.Type)
			unresolvedMap[key]++
			s.UnresolvedCount++
		}

		s.TotalDurationMin += dur

		// Monthly & Weekday
		m := d.Month()
		wd := d.Weekday()

		mm := s.MonthlyStats[m]
		mm.Views++
		mm.DurationMin += dur
		s.MonthlyStats[m] = mm

		wm := s.WeekdayStats[wd]
		wm.Views++
		wm.DurationMin += dur
		s.WeekdayStats[wd] = wm

		// Genre
		for _, g := range genres {
			if _, ok := genreMap[g]; !ok {
				genreMap[g] = &Metric{}
				genreMonthMap[g] = make(map[time.Month]int)
			}
			genreMap[g].Views++
			genreMap[g].DurationMin += dur
			genreMonthMap[g][m] += dur

			// Collect Movie Samples
			if it.Normalized.Type == "movie" {
				if genreSeen[g] == nil {
					genreSeen[g] = make(map[string]bool)
				}
				if !genreSeen[g][it.Normalized.WorkTitle] {
					genreSeen[g][it.Normalized.WorkTitle] = true

					ts := TitleStat{
						Title:       it.Normalized.WorkTitle,
						Type:        it.Normalized.Type,
						DurationMin: dur,
					}
					if it.Metadata != nil {
						ts.VoteAverage = it.Metadata.VoteAverage
						ts.Popularity = it.Metadata.Popularity
						ts.PosterPath = it.Metadata.PosterPath
					}
					s.GenreSampleMovies[g] = append(s.GenreSampleMovies[g], ts)
				}
			}
		}

		// Title
		tKey := fmt.Sprintf("%s|%s", it.Normalized.WorkTitle, it.Normalized.Type) // Use WorkTitle as key
		if _, ok := titleMap[tKey]; !ok {
			titleMap[tKey] = &TitleStat{
				Title: it.Normalized.WorkTitle,
				Type:  it.Normalized.Type,
			}
		}
		titleMap[tKey].Views++
		titleMap[tKey].DurationMin += dur
		if it.Metadata != nil {
			titleMap[tKey].VoteAverage = it.Metadata.VoteAverage
			titleMap[tKey].Popularity = it.Metadata.Popularity
			titleMap[tKey].PosterPath = it.Metadata.PosterPath
			titleGenres[tKey] = it.Metadata.Genres
		}

		// Series
		if it.Normalized.Type == "tv" {
			sn := it.Normalized.WorkTitle // Assuming WorkTitle is Series Name for TV
			if _, ok := seriesMap[sn]; !ok {
				seriesMap[sn] = &SeriesStat{SeriesName: sn, SpanStart: d, SpanEnd: d}
			}
			st := seriesMap[sn]
			st.Views++
			st.DurationMin += dur
			if d.Before(st.SpanStart) {
				st.SpanStart = d
			}
			if d.After(st.SpanEnd) {
				st.SpanEnd = d
			}
		}
	}

	// Calculate Active Days & Ratio
	activeDaysMap := make(map[string]bool)
	for _, d := range dates {
		activeDaysMap[d.Format("2006-01-02")] = true
	}
	s.ActiveDays = len(activeDaysMap)

	// Post-aggregation processing

	// Streaks & Gaps
	s.computeStreaksAndGaps(dates)

	// Genres
	s.computeGenres(genreMap, genreMonthMap)
	s.computeGenreRankings(titleMap, titleGenres)
	s.computeSampleMovies()

	// Titles
	s.computeTitles(titleMap)

	// Series
	s.computeSeries(seriesMap)

	// Unresolved
	s.computeUnresolved(unresolvedMap)

	return s
}

func (s *Stats) computeStreaksAndGaps(dates []time.Time) {
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

	// Check Gaps and Streaks
	var streaks []Streak
	var gaps []Gap

	if len(uniq) > 0 {
		curStart := uniq[0]
		streakLen := 1

		for i := 1; i < len(uniq); i++ {
			prev := uniq[i-1]
			cur := uniq[i]
			diff := cur.Sub(prev) // Should be multiples of 24h approximately
			daysDiff := int(diff.Hours() / 24)

			if daysDiff == 1 {
				streakLen++
			} else {
				// End of streak
				streaks = append(streaks, Streak{Days: streakLen, Start: curStart, End: prev})

				// Found a gap (daysDiff > 1)
				// Gap is between prev and cur.
				// Gap length is daysDiff - 1 (days without view)
				// e.g. 1st and 3rd viewed. diff=2. gap=1 (2nd).
				gapStart := prev.Add(24 * time.Hour)
				gapEnd := cur.Add(-24 * time.Hour)
				gaps = append(gaps, Gap{Days: daysDiff - 1, Start: gapStart, End: gapEnd})

				streakLen = 1
				curStart = cur
			}
		}
		// Final streak
		streaks = append(streaks, Streak{Days: streakLen, Start: curStart, End: uniq[len(uniq)-1]})
	}

	// Sort Streaks
	sort.Slice(streaks, func(i, j int) bool {
		return streaks[i].Days > streaks[j].Days
	})
	if len(streaks) > 3 {
		s.TopStreaks = streaks[:3]
	} else {
		s.TopStreaks = streaks
	}

	// Max Gap
	for _, g := range gaps {
		if g.Days > s.MaxGap.Days {
			s.MaxGap = g
		}
	}
}

func (s *Stats) computeGenres(m map[string]*Metric, mm map[string]map[time.Month]int) {
	var gs []GenreStat
	for name, met := range m {
		gs = append(gs, GenreStat{
			Name:        name,
			DurationMin: met.DurationMin,
			Views:       met.Views,
			Share:       0, // filled later
		})
	}

	// Sort by Duration Desc
	sort.Slice(gs, func(i, j int) bool {
		return gs[i].DurationMin > gs[j].DurationMin
	})

	// Fill Share
	if s.TotalDurationMin > 0 {
		for i := range gs {
			gs[i].Share = float64(gs[i].DurationMin) / float64(s.TotalDurationMin) * 100
		}
	}

	s.GenreStats = gs

	// Compute Spikes
	s.GenreMonthSpike = make(map[string]Spike)
	for name, monthDur := range mm {
		var maxM time.Month
		var maxD int
		for m, d := range monthDur {
			if d > maxD {
				maxD = d
				maxM = m
			}
		}
		if maxD > 0 {
			s.GenreMonthSpike[name] = Spike{Month: maxM, DurationMin: maxD}
		}
	}
}

func (s *Stats) computeTitles(m map[string]*TitleStat) {
	var ts []TitleStat
	for _, v := range m {
		ts = append(ts, *v)
	}

	// By Duration
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].DurationMin > ts[j].DurationMin
	})
	if len(ts) > 10 {
		s.TopTitlesByDuration = ts[:10] // limit to top 10 for now, or as needed
	} else {
		s.TopTitlesByDuration = ts
	}
	// Or maybe user wants top 30? Plan said "Top" without number, or Top 30 for unresolved.
	// But let's keep all or top 20. Let's start with all, and filter in render, OR stick to plan hint.
	// Actually for "Top Titles", template implies just top few. Let's store enough (e.g. 50) and slice in render.
	// Wait, Slice is destructive if I re-sort. Copying.

	ts2 := make([]TitleStat, len(ts))
	copy(ts2, ts)

	// By Views
	sort.Slice(ts2, func(i, j int) bool {
		return ts2[i].Views > ts2[j].Views
	})

	// actually let's just store top 50 strictly
	limit := 50
	if len(ts) > limit {
		s.TopTitlesByDuration = ts[:limit]
	} else {
		s.TopTitlesByDuration = ts
	}

	if len(ts2) > limit {
		s.TopTitlesByViews = ts2[:limit]
	} else {
		s.TopTitlesByViews = ts2
	}
}

func (s *Stats) computeSeries(m map[string]*SeriesStat) {
	var ss []SeriesStat
	for _, v := range m {
		ss = append(ss, *v)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].DurationMin > ss[j].DurationMin
	})

	limit := 50
	if len(ss) > limit {
		s.TopSeriesByDuration = ss[:limit]
	} else {
		s.TopSeriesByDuration = ss
	}

	// If need by views
	ss2 := make([]SeriesStat, len(ss))
	copy(ss2, ss)
	sort.Slice(ss2, func(i, j int) bool {
		return ss2[i].Views > ss2[j].Views
	})
	if len(ss2) > limit {
		s.TopSeriesByViews = ss2[:limit]
	} else {
		s.TopSeriesByViews = ss2
	}
}

func (s *Stats) computeUnresolved(m map[string]int) {
	var us []UnresolvedItem
	for k, count := range m {
		// k is "Title|Type"
		var title, typ string
		// simple split might fail if title contains |, but we used fmt.Sprintf earlier.
		// Actually let's just do manual split or parsing.
		// Since we control key gen, let's just parsing.
		// Wait, simple split is better.
		// Or assume no pipe in title? Title comes from Normalized which is safe?
		// Normalized WorkTitle *might* have pipe.
		// Let's rely on range loop earlier storing it properly? No, stats func is monolithic.
		// Let's re-parse key or better: store struct in map?
		// Too late for that in this function scope without changin previous scope.
		// Recalculating map key:
		// Let's use string split from end or just fix valid separator that can't be in values?
		// For now simple split on last | is safest if type is fixed enum.

		// Actually, let's just use string manipulation
		lastPipe := -1
		for i := len(k) - 1; i >= 0; i-- {
			if k[i] == '|' {
				lastPipe = i
				break
			}
		}
		if lastPipe != -1 {
			title = k[:lastPipe]
			typ = k[lastPipe+1:]
		} else {
			title = k
			typ = "unknown"
		}

		us = append(us, UnresolvedItem{
			Title: title,
			Type:  typ,
			Views: count,
		})
	}

	sort.Slice(us, func(i, j int) bool {
		return us[i].Views > us[j].Views
	})

	if len(us) > 30 {
		s.UnresolvedList = us[:30]
	} else {
		s.UnresolvedList = us
	}
}

func (s *Stats) computeGenreRankings(titleMap map[string]*TitleStat, titleGenres map[string][]string) {
	// 1. Group titles by genre
	genreWorks := make(map[string][]TitleStat)
	for tKey, genres := range titleGenres {
		stat, ok := titleMap[tKey]
		if !ok {
			continue
		}
		// Only consider movies/tv that have a vote average > 0 to avoid "unrated" being worst
		// And maybe only "movie"? Request said "movies" (Top 3 Works).
		// "Top 3 Works" (Sakuhin) implies both.
		// If 0 vote average, usually means unrated. Let's keep them out of ranking.
		if stat.VoteAverage <= 0 {
			continue
		}

		for _, g := range genres {
			genreWorks[g] = append(genreWorks[g], *stat)
		}
	}

	for g, works := range genreWorks {
		// Sort for Top (Desc)
		sort.Slice(works, func(i, j int) bool {
			return works[i].VoteAverage > works[j].VoteAverage
		})

		// Top 3
		topLimit := 3
		if len(works) > topLimit {
			s.GenreTopWorks[g] = works[:topLimit]
		} else {
			s.GenreTopWorks[g] = works
		}

		// Sort for Worst (Asc) -> Filtered > 0 already
		// Need to copy relevant slice first?
		// works is already sorted desc.
		// Last elements are the worst.
		// But slice might be large.
		// Let's just create a new slice for worst or just pick from end if len is enough.
		// But we need exactly 3 worst.
		// If len < 3, they are both top and worst? Overlap is fine.

		// Let's take the last 3 elements, reverse them so it's "Worst 1, Worst 2, Worst 3" (lowest first)
		worst := make([]TitleStat, 0, 3)
		n := len(works)
		if n > 0 {
			// End is lowest.
			// works[n-1] is lowest. works[n-2] is second lowest.
			// We want: [Lowest, 2nd Lowest, 3rd Lowest]
			for i := n - 1; i >= 0 && len(worst) < 3; i-- {
				worst = append(worst, works[i])
			}
		}
		s.GenreWorstWorks[g] = worst
	}
}

func (s *Stats) computeSampleMovies() {
	for g, movies := range s.GenreSampleMovies {
		// Sort by Popularity Desc
		sort.Slice(movies, func(i, j int) bool {
			return movies[i].Popularity > movies[j].Popularity
		})
		// Truncate to 5
		if len(movies) > 5 {
			s.GenreSampleMovies[g] = movies[:5]
		}
	}
}
