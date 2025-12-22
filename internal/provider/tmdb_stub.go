package provider

import (
	"strings"

	"github.com/kmdkuk/nfrecap/internal/model"
)

type TMDbStub struct{}

func NewTMDbStub() *TMDbStub { return &TMDbStub{} }

// NOTE: これは雛形用スタブ。後でTMDb実装に差し替える。
func (p *TMDbStub) Lookup(workTitle string, typ string) (model.Metadata, bool, error) {
	// 超雑：とりあえずそれっぽいメタデータを返す
	md := model.Metadata{
		Provider: "tmdb-stub",
		ID:       "stub:" + strings.ToLower(workTitle),
		Title:    workTitle,
	}
	if typ == "movie" {
		md.Genres = []string{"Unknown"}
		md.Runtime = 120
	} else if typ == "tv" {
		md.Genres = []string{"Unknown"}
		md.Runtime = 45
	}
	return md, true, nil
}
