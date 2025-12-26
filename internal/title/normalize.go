package title

import (
	"strings"

	"github.com/kmdkuk/nfrecap/internal/model"
)

func Normalize(raw string) model.NormalizedTitle {
	s := strings.TrimSpace(raw)
	n := model.NormalizedTitle{
		RawTitle:  s,
		WorkTitle: s,
		Type:      "movie",
	}

	// Common Netflix JP format: "作品名: シーズン1: 話タイトル"
	parts := strings.Split(s, ":")
	if len(parts) >= 1 {
		n.WorkTitle = strings.TrimSpace(parts[0])
	}

	if len(parts) >= 2 {
		n.Type = "tv"
		n.Season = strings.TrimSpace(parts[1])
	}
	if len(parts) >= 3 {
		n.EpisodeTitle = strings.TrimSpace(parts[2])
	}

	return n
}
