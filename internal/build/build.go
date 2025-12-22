package build

import (
	"encoding/json"
	"time"

	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/kmdkuk/nfrecap/internal/provider"
	"github.com/kmdkuk/nfrecap/internal/store"
	"github.com/kmdkuk/nfrecap/internal/title"
)

type Options struct {
	Fetch   bool
	Verbose bool
}

type Summary struct {
	CacheHits   int
	CacheMisses int
	Fetched     int
	Unresolved  int
}

type Built struct {
	Source      string      `json:"source,omitempty"`
	GeneratedAt string      `json:"generated_at"`
	Items       []BuiltItem `json:"items"`
}

type BuiltItem struct {
	Date       string                `json:"date"`
	Normalized model.NormalizedTitle `json:"normalized"`
	Metadata   *model.Metadata       `json:"metadata,omitempty"`
}

func Run(records []model.ViewingRecord, cache store.Cache, p provider.Provider, opts Options) ([]byte, Summary, error) {
	sum := Summary{}
	out := Built{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Items:       make([]BuiltItem, 0, len(records)),
	}

	for _, r := range records {
		n := title.Normalize(r.Title)

		md, ok, err := cache.Get(n.WorkTitle, n.Type)
		if err != nil {
			return nil, sum, err
		}
		if ok {
			sum.CacheHits++
			out.Items = append(out.Items, BuiltItem{
				Date:       r.Date.Format("2006-01-02"),
				Normalized: n,
				Metadata:   &md,
			})
			continue
		}

		sum.CacheMisses++

		if opts.Fetch {
			got, found, err := p.Lookup(n.WorkTitle, n.Type)
			if err != nil {
				return nil, sum, err
			}
			if found {
				sum.Fetched++
				_ = cache.Put(n.WorkTitle, n.Type, got)

				out.Items = append(out.Items, BuiltItem{
					Date:       r.Date.Format("2006-01-02"),
					Normalized: n,
					Metadata:   &got,
				})
				continue
			}
		}

		sum.Unresolved++
		out.Items = append(out.Items, BuiltItem{
			Date:       r.Date.Format("2006-01-02"),
			Normalized: n,
			Metadata:   nil,
		})
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, sum, err
	}
	return b, sum, nil
}
