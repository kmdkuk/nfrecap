package recap

import (
	"fmt"
	"strings"
	"time"
)

func RenderMarkdown(s Stats) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# Netflix Recap %d\n\n", s.Year)

	fmt.Fprintf(&b, "## 視聴本数\n- 合計: %d\n\n", s.TotalItems)

	fmt.Fprintf(&b, "## 月別視聴本数\n")
	for i := 0; i < 12; i++ {
		fmt.Fprintf(&b, "- %02d月: %d\n", i+1, s.ItemsByMonth[i])
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "## 曜日別視聴本数\n")
	wd := []string{"日", "月", "火", "水", "木", "金", "土"}
	for i := 0; i < 7; i++ {
		fmt.Fprintf(&b, "- %s: %d\n", wd[i], s.ItemsByWeekday[i])
	}
	b.WriteString("\n")

	fmt.Fprintf(&b, "## 最長連続視聴（streak）\n")
	if s.LongestStreak <= 1 || s.LongestStreakStart == nil || s.LongestStreakEnd == nil {
		fmt.Fprintf(&b, "- 最長: %d日\n\n", s.LongestStreak)
	} else {
		fmt.Fprintf(&b, "- 最長: %d日\n", s.LongestStreak)
		fmt.Fprintf(&b, "- 期間: %s 〜 %s\n\n", formatDate(*s.LongestStreakStart), formatDate(*s.LongestStreakEnd))
	}

	return b.String()
}

func formatDate(t time.Time) string { return t.Format("2006-01-02") }
