package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kmdkuk/nfrecap/internal/build"
	"github.com/kmdkuk/nfrecap/internal/csvio"
	"github.com/kmdkuk/nfrecap/internal/provider"
	"github.com/kmdkuk/nfrecap/internal/store"
)

var (
	buildIn       string
	buildOut      string
	buildFetch    bool
	buildCacheDir string
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

		cache := store.NewFileCache(buildCacheDir)

		// Providerはスタブ（後でTMDb実装に置き換え）
		p := provider.NewTMDbStub()

		opts := build.Options{
			Fetch:   buildFetch,
			Verbose: flagVerbose,
		}

		out, summary, err := build.Run(recs, cache, p, opts)
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

	_ = buildCmd.MarkFlagRequired("in")
}
