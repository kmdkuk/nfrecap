package build

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"

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
		Items:       make([]BuiltItem, len(records)),
	}

	var mu sync.Mutex
	eg, ctx := errgroup.WithContext(context.Background())
	// Limit: 40 req/sec, Burst: 1
	limiter := rate.NewLimiter(rate.Limit(40), 1)

	for i, r := range records {
		i, r := i, r // capture loop variables
		eg.Go(func() error {
			n := title.Normalize(r.Title)

			// Cache Read (RLock-like behavior, but using Mutex for simplicity across all cache ops)
			mu.Lock()
			md, ok, err := cache.Get(n.WorkTitle, n.Type)
			mu.Unlock()

			if err != nil {
				return err
			}

			if ok {
				mu.Lock()
				sum.CacheHits++
				mu.Unlock()

				out.Items[i] = BuiltItem{
					Date:       r.Date.Format("2006-01-02"),
					Normalized: n,
					Metadata:   &md,
				}
				return nil
			}

			mu.Lock()
			sum.CacheMisses++
			mu.Unlock()

			if opts.Fetch {
				// Rate Limit
				if err := limiter.Wait(ctx); err != nil {
					return err
				}

				got, found, err := p.Lookup(n.WorkTitle, n.Type)
				if err != nil {
					return err
				}
				if found {
					mu.Lock()
					sum.Fetched++
					// Cache Write
					putErr := cache.Put(n.WorkTitle, n.Type, got)
					mu.Unlock()

					if putErr != nil {
						// caching error shouldn't stop the build? 
						// but current logic returns error. keeping consistent.
						return putErr
					}

					out.Items[i] = BuiltItem{
						Date:       r.Date.Format("2006-01-02"),
						Normalized: n,
						Metadata:   &got,
					}
					return nil
				}
			}

			mu.Lock()
			sum.Unresolved++
			mu.Unlock()

			out.Items[i] = BuiltItem{
				Date:       r.Date.Format("2006-01-02"),
				Normalized: n,
				Metadata:   nil,
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, sum, err
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, sum, err
	}
	return b, sum, nil
}
