package recap

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"
)

func RenderMarkdown(s Stats) string {
	tmplStr := `
# Netflix Recap {{.Year}}

> 生成日時: {{.GeneratedAt}}

---

## まずは要点（TL;DR）

- 推定総視聴時間：**{{.TotalDurationHours}} 時間**（{{.TotalDurationMin}} 分）
- 視聴回数：**{{.TotalViews}} 本**
- 視聴日数：**{{.ActiveDays}} 日**（{{.ActiveRatio}}%）
- 最長連続視聴：**{{.LongestStreakDays}} 日**（{{.LongestStreakStart}} 〜 {{.LongestStreakEnd}}）
- メタデータ取得率：**{{.CoveredViews}} / {{.TotalViews}}**（{{.CoverageRatio}}%）

---

## 1. 全体概要

### 月別 推定視聴時間・視聴回数

| 月 | 視聴回数 | 推定視聴時間（時間） |
|---:|---:|---:|
{{- range .MonthlyRows }}
| {{.Month}} | {{.Views}} | {{.Hours}} |
{{- end }}

---

### 曜日別 視聴傾向

| 曜日 | 視聴回数 | 推定視聴時間（時間） |
|---|---:|---:|
{{- range .WeekdayRows }}
| {{.Weekday}} | {{.Views}} | {{.Hours}} |
{{- end }}

> ※ 推定視聴時間は ` + "`runtime_min`" + ` を基に算出しています
> ※ TV シリーズの場合、1話あたりの代表的な再生時間を使用しています

---

## 2. 視聴の継続性（Streak）

### 最長連続視聴記録

- **{{.LongestStreakDays}} 日連続**
- 期間：{{.LongestStreakStart}} 〜 {{.LongestStreakEnd}}

---

### 連続視聴ランキング（Top 3）

| 順位 | 連続日数 | 期間 |
|---:|---:|---|
{{- range .TopStreaksRows }}
| {{.Rank}}位 | {{.Days}} 日 | {{.Start}} 〜 {{.End}} |
{{- end }}

---

### 視聴活動の間隔

- 視聴があった日数：{{.ActiveDays}} 日
- 最長の空白期間：{{.MaxGapDays}} 日
  （{{.MaxGapStart}} 〜 {{.MaxGapEnd}}）

---

## 3. ジャンル別分析（時間ベース）

### ジャンル別 推定視聴時間（Top 10 + その他）

| ジャンル | 推定視聴時間（時間） | 割合 | 視聴回数 |
|---|---:|---:|---:|
{{- range .GenreRows }}
| {{.Name}} | {{.Hours}} | {{.Share}}% | {{.Views}} |
{{- end }}

---

### ジャンルの偏り（月別ピーク）

| ジャンル | 最多視聴月 | ピーク時間（時間） | 備考 |
|---|---|---:|---|
{{- range .GenreSpikeRows }}
| {{.Name}} | {{.Month}} | {{.Hours}} | {{.Note}} |
{{- end }}

---

## 4. 作品・シリーズ別


---


### ジャンル別視聴映画（サンプル）

| ジャンル | 代表的な視聴作品（3選） |
|---|---|
{{- range .GenreRows }}
{{- if .SampleMovies }}
| {{.Name}} | {{.SampleMovies}} |
{{- end }}
{{- end }}

---

### 推定視聴時間が多い作品（Top）

| 順位 | 作品名 | 種別 | 推定視聴時間（時間） | 視聴回数 |
|---:|---|---|---:|---:|
{{- range .TopTitlesByDurationRows }}
| {{.Rank}}位 | {{.Title}} | {{.Type}} | {{.Hours}} | {{.Views}} |
{{- end }}


---

### シリーズ視聴（TV作品のみ）

| 順位 | シリーズ名 | 視聴回数 | 推定視聴時間（時間） | 集中視聴期間 |
|---:|---|---:|---:|---|
{{- range .TopSeriesRows }}
| {{.Rank}}位 | {{.SeriesName}} | {{.Views}} | {{.Hours}} | {{.Span}} |
{{- end }}

---

## 5. データ品質・補足

### メタデータ取得状況

- 取得済み：**{{.CoveredViews}} / {{.TotalViews}}**（{{.CoverageRatio}}%）
- 未取得（unresolved）：{{.UnresolvedCount}} 件

---

### 注記

- 本文章は [https://github.com/kmdkuk/nfrecap](https://github.com/kmdkuk/nfrecap) を利用して生成されました。
- 作品の特定およびメタデータ取得には [https://www.themoviedb.org/](https://www.themoviedb.org/) を利用しています。
- Disclaimer: This nfrecap uses TMDB and the TMDB APIs but is not endorsed, certified, or otherwise approved by TMDB.
- 推定視聴時間は参考値であり、実際の再生時間と一致しない場合があります。
- 未取得作品は、今後の正規化ルール改善や手動補正で解消できる可能性があります。
`

	data := prepareViewData(s)

	t, err := template.New("recap").Parse(tmplStr)
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return buf.String()
}

// view data structs
type viewData struct {
	Year               int
	GeneratedAt        string
	SourceFile         string
	TotalDurationHours string
	TotalDurationMin   int
	TotalViews         int
	ActiveDays         int
	ActiveRatio        string
	LongestStreakDays  int
	LongestStreakStart string
	LongestStreakEnd   string
	CoveredViews       int
	CoverageRatio      string
	UnresolvedCount    int
	MaxGapDays         int
	MaxGapStart        string
	MaxGapEnd          string

	MonthlyRows             []monthlyRow
	WeekdayRows             []weekdayRow
	TopStreaksRows          []streakRow
	GenreRows               []genreRow
	GenreSpikeRows          []spikeRow
	TopTitlesByDurationRows []titleRow
	TopTitlesByViewsRows    []titleRow
	TopSeriesRows           []seriesRow
	UnresolvedRows          []unresolvedRow
}

type monthlyRow struct {
	Month   string
	Views   int
	Hours   string
	Minutes int
}
type weekdayRow struct {
	Weekday string
	Views   int
	Hours   string
}
type streakRow struct {
	Rank  int
	Days  int
	Start string
	End   string
}
type genreRow struct {
	Name         string
	Hours        string
	Share        string
	Views        int
	SampleMovies string // comma separated or pre-formatted? Or separate section.
}
type spikeRow struct {
	Name  string
	Month int
	Hours string
	Note  string
}
type titleRow struct {
	Rank  int
	Title string
	Type  string
	Hours string
	Views int
}
type seriesRow struct {
	Rank       int
	SeriesName string
	Views      int
	Hours      string
	Span       string
}
type unresolvedRow struct {
	Rank  int
	Title string
	Type  string
	Views int
}

func prepareViewData(s Stats) viewData {
	vd := viewData{
		Year:             s.Year,
		GeneratedAt:      s.GeneratedAt,
		SourceFile:       s.SourceFile,
		TotalDurationMin: s.TotalDurationMin,
		TotalViews:       s.TotalViews,
		ActiveDays:       s.ActiveDays,
		UnresolvedCount:  s.UnresolvedCount,
	}

	vd.TotalDurationHours = fmt.Sprintf("%.1f", float64(s.TotalDurationMin)/60.0)
	vd.ActiveRatio = fmt.Sprintf("%.1f", float64(s.ActiveDays)/365.0*100.0) // simplified 365

	covered := s.TotalViews - s.UnresolvedCount
	vd.CoveredViews = covered
	if s.TotalViews > 0 {
		vd.CoverageRatio = fmt.Sprintf("%.1f", float64(covered)/float64(s.TotalViews)*100.0)
	} else {
		vd.CoverageRatio = "0.0"
	}

	// Longest Streak
	if len(s.TopStreaks) > 0 {
		top := s.TopStreaks[0]
		vd.LongestStreakDays = top.Days
		vd.LongestStreakStart = top.Start.Format("2006-01-02")
		vd.LongestStreakEnd = top.End.Format("2006-01-02")
	}

	// Max Gap
	vd.MaxGapDays = s.MaxGap.Days
	if s.MaxGap.Days > 0 {
		vd.MaxGapStart = s.MaxGap.Start.Format("2006-01-02")
		vd.MaxGapEnd = s.MaxGap.End.Format("2006-01-02")
	}

	// Monthly
	for m := time.January; m <= time.December; m++ {
		metric := s.MonthlyStats[m]
		vd.MonthlyRows = append(vd.MonthlyRows, monthlyRow{
			Month:   fmt.Sprintf("%d月", m),
			Views:   metric.Views,
			Hours:   fmt.Sprintf("%.1f", float64(metric.DurationMin)/60.0),
			Minutes: metric.DurationMin,
		})
	}

	// Weekday
	weekdays := []string{"日", "月", "火", "水", "木", "金", "土"}
	for i, w := range weekdays {
		metric := s.WeekdayStats[time.Weekday(i)]
		vd.WeekdayRows = append(vd.WeekdayRows, weekdayRow{
			Weekday: w,
			Views:   metric.Views,
			Hours:   fmt.Sprintf("%.1f", float64(metric.DurationMin)/60.0),
		})
	}

	// Top Streaks
	for i, st := range s.TopStreaks {
		vd.TopStreaksRows = append(vd.TopStreaksRows, streakRow{
			Rank:  i + 1,
			Days:  st.Days,
			Start: st.Start.Format("2006-01-02"),
			End:   st.End.Format("2006-01-02"),
		})
	}

	// Genres
	// Limit top 10 + others
	otherGenreDur := 0
	otherGenreViews := 0

	limit := 10
	usedMovies := make(map[string]bool) // Track movies already used in samples

	for i, g := range s.GenreStats {
		// Prepare samples string with deduplication
		samples := ""
		if movies, ok := s.GenreSampleMovies[g.Name]; ok && len(movies) > 0 {
			// Filter out already used movies
			var availableTitles []string
			for _, movie := range movies {
				if !usedMovies[movie.Title] {
					availableTitles = append(availableTitles, movie.Title)
				}
			}

			// Take up to 3 unique movies
			sampleLimit := 3
			if len(availableTitles) < sampleLimit {
				sampleLimit = len(availableTitles)
			}

			if sampleLimit > 0 {
				selectedTitles := availableTitles[:sampleLimit]
				samples = strings.Join(selectedTitles, ", ")

				// Mark these movies as used
				for _, title := range selectedTitles {
					usedMovies[title] = true
				}
			}
		}

		if i < limit {
			vd.GenreRows = append(vd.GenreRows, genreRow{
				Name:         g.Name,
				Hours:        fmt.Sprintf("%.1f", float64(g.DurationMin)/60.0),
				Share:        fmt.Sprintf("%.1f", g.Share),
				Views:        g.Views,
				SampleMovies: samples,
			})
		} else {
			otherGenreDur += g.DurationMin
			otherGenreViews += g.Views
		}
	}
	if otherGenreDur > 0 {
		share := 0.0
		if s.TotalDurationMin > 0 {
			share = float64(otherGenreDur) / float64(s.TotalDurationMin) * 100.0
		}
		vd.GenreRows = append(vd.GenreRows, genreRow{
			Name:  "その他",
			Hours: fmt.Sprintf("%.1f", float64(otherGenreDur)/60.0),
			Share: fmt.Sprintf("%.1f", share),
			Views: otherGenreViews,
		})
	}

	// Spikes
	// Just show top genres spikes? Or all? User said "Spike1, ..." table.
	// Let's show spikes for Top 10 genres.
	for i, g := range s.GenreStats {
		if i >= 10 {
			break
		}
		sp, ok := s.GenreMonthSpike[g.Name]
		if ok {
			vd.GenreSpikeRows = append(vd.GenreSpikeRows, spikeRow{
				Name:  g.Name,
				Month: int(sp.Month),
				Hours: fmt.Sprintf("%.1f", float64(sp.DurationMin)/60.0),
				Note:  "", // Placeholder
			})
		}
	}

	// Sort Spikes by Month
	sort.Slice(vd.GenreSpikeRows, func(i, j int) bool {
		return vd.GenreSpikeRows[i].Month < vd.GenreSpikeRows[j].Month
	})

	// Titles
	for i, t := range s.TopTitlesByDurationRows(s.TopTitlesByDuration) {
		vd.TopTitlesByDurationRows = append(vd.TopTitlesByDurationRows, titleRow{
			Rank:  i + 1,
			Title: t.Title,
			Type:  t.Type,
			Hours: fmt.Sprintf("%.1f", float64(t.DurationMin)/60.0),
			Views: t.Views,
		})
	}
	for i, t := range s.TopTitlesByViewsRows(s.TopTitlesByViews) {
		vd.TopTitlesByViewsRows = append(vd.TopTitlesByViewsRows, titleRow{
			Rank:  i + 1,
			Title: t.Title,
			Type:  t.Type,
			Hours: fmt.Sprintf("%.1f", float64(t.DurationMin)/60.0),
			Views: t.Views,
		})
	}

	// Series
	for i, ser := range s.TopSeriesByDuration {
		if i >= 10 {
			break
		}
		vd.TopSeriesRows = append(vd.TopSeriesRows, seriesRow{
			Rank:       i + 1,
			SeriesName: ser.SeriesName,
			Views:      ser.Views,
			Hours:      fmt.Sprintf("%.1f", float64(ser.DurationMin)/60.0),
			Span:       fmt.Sprintf("%s 〜 %s", ser.SpanStart.Format("2006-01-02"), ser.SpanEnd.Format("2006-01-02")),
		})
	}

	// Unresolved
	for i, u := range s.UnresolvedList {
		vd.UnresolvedRows = append(vd.UnresolvedRows, unresolvedRow{
			Rank:  i + 1,
			Title: u.Title,
			Type:  u.Type,
			Views: u.Views,
		})
	}

	return vd
}

// Helpers for slice conversion
func (s *Stats) TopTitlesByDurationRows(ts []TitleStat) []TitleStat {
	if len(ts) > 10 {
		return ts[:10]
	}
	return ts
}
func (s *Stats) TopTitlesByViewsRows(ts []TitleStat) []TitleStat {
	if len(ts) > 10 {
		return ts[:10]
	}
	return ts
}
