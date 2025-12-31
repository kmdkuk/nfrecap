package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/kmdkuk/nfrecap/internal/build"
	"github.com/kmdkuk/nfrecap/internal/csvio"
	tmdbprovider "github.com/kmdkuk/nfrecap/internal/provider/tmdb"
	"github.com/kmdkuk/nfrecap/internal/store"
)

var (
	buildIn       string
	buildOut      string
	buildFetch    bool
	buildCacheDir string
	buildCacheTTL time.Duration
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build normalized recap data from Netflix viewing history",
	Long: `Build constructs normalized recap data from Netflix viewing history CSV.

By default, it uses locally cached metadata only (no network).
Use --fetch to retrieve metadata from external APIs and update the cache.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		recs, err := csvio.ReadNetflixCSV(buildIn)
		if err != nil {
			return err
		}

		cache := store.NewFileCache(buildCacheDir, buildCacheTTL)

		// Providerの初期化 (TMDB API)
		p, err := tmdbprovider.NewFromEnv(tmdbprovider.Options{
			UseV4Bearer: true,
			AutoRetry:   true,
			Language:    "ja-JP", // TODO: Configurable
		})
		if err != nil {
			return fmt.Errorf("failed to init tmdb provider: %w", err)
		}

		opts := build.Options{
			Fetch:   buildFetch,
			Verbose: flagVerbose,
		}

		outStruct, summary, err := build.Run(recs, cache, p, opts)
		if err != nil {
			return err
		}

		out, err := json.MarshalIndent(outStruct, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(buildOut, out, 0644); err != nil {
			return err
		}

		if flagVerbose {
			fmt.Fprintf(os.Stderr, "wrote %s\n", buildOut)
			fmt.Fprintf(os.Stderr, "cache hits=%d misses=%d fetched=%d unresolved=%d\n",
				summary.CacheHits, summary.CacheMisses, summary.Fetched, summary.Unresolved)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&buildIn, "in", "i", "", "input Netflix CSV file")
	buildCmd.Flags().StringVarP(&buildOut, "out", "o", "build.json", "output built JSON file")
	buildCmd.Flags().BoolVar(&buildFetch, "fetch", false, "fetch metadata from external APIs before building")
	buildCmd.Flags().StringVar(&buildCacheDir, "cache-dir", store.DefaultCacheDir(), "metadata cache directory")
	buildCmd.Flags().DurationVar(&buildCacheTTL, "cache-ttl", 72*time.Hour, "cache expiration duration")

	_ = buildCmd.MarkFlagRequired("in")
}
