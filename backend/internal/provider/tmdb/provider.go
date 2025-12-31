package tmdbprovider

import (
	"errors"
	"fmt"
	"os"
	"strings"

	tmdb "github.com/cyruzin/golang-tmdb"

	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/kmdkuk/nfrecap/internal/provider"
)

// Ensure implements provider.Provider
var _ provider.Provider = (*Provider)(nil)

type Provider struct {
	c *tmdb.Client
}

type Options struct {
	UseV4Bearer bool // true: InitV4, false: Init (api key)
	AutoRetry   bool // retry on 429
	Language    string // e.g. "ja-JP"
}

func NewFromEnv(opts Options) (*Provider, error) {
	var (
		client *tmdb.Client
		err    error
	)

	// TODO: Bearer token get from viper
	if opts.UseV4Bearer {
		token := strings.TrimSpace(os.Getenv("TMDB_BEARER_TOKEN"))
		if token == "" {
			return nil, errors.New("TMDB_BEARER_TOKEN is not set")
		}
		client, err = tmdb.InitV4(token) // Bearer token init :contentReference[oaicite:4]{index=4}
		if err != nil {
			return nil, err
		}
	} else {
		key := strings.TrimSpace(os.Getenv("TMDB_API_KEY"))
		if key == "" {
			return nil, errors.New("TMDB_API_KEY is not set")
		}
		client, err = tmdb.Init(key) // API key init :contentReference[oaicite:5]{index=5}
		if err != nil {
			return nil, err
		}
	}

	if opts.AutoRetry {
		client.SetClientAutoRetry() // 429 retry helper :contentReference[oaicite:6]{index=6}
	}

	return &Provider{c: client}, nil
}

func (p *Provider) Lookup(workTitle string, typ string) (model.Metadata, bool, error) {
	var (
		id   int64
		kind string
	)

	// 1. Search
	switch typ {
	case "movie":
		res, err := p.c.GetSearchMovies(workTitle, nil)
		if err != nil {
			return model.Metadata{}, false, err
		}
		if len(res.Results) == 0 {
			return model.Metadata{}, false, nil
		}
		id = res.Results[0].ID
		kind = "movie"
	case "tv":
		res, err := p.c.GetSearchTVShow(workTitle, nil)
		if err != nil {
			return model.Metadata{}, false, err
		}
		if len(res.Results) == 0 {
			return model.Metadata{}, false, nil
		}
		id = res.Results[0].ID
		kind = "tv"
	default:
		// Unknown type: try movie first, then tv
		resM, err := p.c.GetSearchMovies(workTitle, nil)
		if err == nil && len(resM.Results) > 0 {
			id = resM.Results[0].ID
			kind = "movie"
		} else {
			resTV, err := p.c.GetSearchTVShow(workTitle, nil)
			if err == nil && len(resTV.Results) > 0 {
				id = resTV.Results[0].ID
				kind = "tv"
			} else {
				// both failed or found nothing
				return model.Metadata{}, false, nil
			}
		}
	}

	// 2. Get Details
	if kind == "movie" {
		details, err := p.c.GetMovieDetails(int(id), nil)
		if err != nil {
			return model.Metadata{}, false, err
		}

		genres := make([]string, len(details.Genres))
		for i, g := range details.Genres {
			genres[i] = g.Name
		}

		return model.Metadata{
			Provider: "tmdb",
			ID:       fmt.Sprintf("movie:%d", id),
			Title:    details.Title,
			Genres:   genres,
			Runtime:  int(details.Runtime),
		}, true, nil

	} else { // tv
		details, err := p.c.GetTVDetails(int(id), nil)
		if err != nil {
			return model.Metadata{}, false, err
		}

		genres := make([]string, len(details.Genres))
		for i, g := range details.Genres {
			genres[i] = g.Name
		}

		// TV usually has "episode_run_time" (array) or we can just not set Runtime if ambiguous.
		// Taking the first value if available, or just leave it 0.
		runtime := 0
		if len(details.EpisodeRunTime) > 0 {
			runtime = int(details.EpisodeRunTime[0])
		}

		return model.Metadata{
			Provider: "tmdb",
			ID:       fmt.Sprintf("tv:%d", id),
			Title:    details.Name,
			Genres:   genres,
			Runtime:  runtime,
		}, true, nil
	}
}
