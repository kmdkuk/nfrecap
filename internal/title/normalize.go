package title

import (
	"regexp"
	"strings"

	"github.com/kmdkuk/nfrecap/internal/model"
)

var reSeason = regexp.MustCompile(`シーズン\s*([0-9]+)`)

func Normalize(raw string) model.NormalizedTitle {
	s := strings.TrimSpace(raw)
	n := model.NormalizedTitle{
		RawTitle:  s,
		WorkTitle: s,
		Type:      "movie",
	}

	// Heuristic: if contains "シーズン" treat as TV
	if strings.Contains(s, "シーズン") {
		n.Type = "tv"

		// Common Netflix JP format: "作品名: シーズン1: 話タイトル"
		parts := strings.Split(s, ":")
		if len(parts) >= 1 {
			n.WorkTitle = strings.TrimSpace(parts[0])
		}

		if m := reSeason.FindStringSubmatch(s); len(m) == 2 {
			n.Season = atoiSafe(m[1])
		}
		if len(parts) >= 3 {
			n.EpisodeTitle = strings.TrimSpace(parts[2])
		}
	}

	return n
}

func atoiSafe(s string) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			continue
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
