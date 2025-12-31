package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/kmdkuk/nfrecap/internal/build"
	"github.com/kmdkuk/nfrecap/internal/csvio"
	tmdbprovider "github.com/kmdkuk/nfrecap/internal/provider/tmdb"
	"github.com/kmdkuk/nfrecap/internal/recap"
	"github.com/kmdkuk/nfrecap/internal/store"
)

var (
	servePort     int
	serveCacheDir string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server to process Netflix CSV and return recap JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Init dependencies
		cache := store.NewFileCache(serveCacheDir)

		p, err := tmdbprovider.NewFromEnv(tmdbprovider.Options{
			UseV4Bearer: true,
			AutoRetry:   true,
			Language:    "ja-JP",
		})
		if err != nil {
			return fmt.Errorf("failed to init tmdb provider: %w", err)
		}

		// Setup Handler
		http.HandleFunc("/api/recap", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Parse Multipart
			// 10MB limit
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			file, _, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Missing file part", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Parse CSV
			records, err := csvio.ParseNetflixCSV(file)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to parse CSV: %v", err), http.StatusBadRequest)
				return
			}

			// Build
			opts := build.Options{
				Fetch:   true, // Always fetch (or make it configurable via query param?)
				Verbose: true, // Log to stdout/stderr
			}

			// Recap
			yearVal := r.FormValue("year")
			year := time.Now().Year()
			if yearVal != "" {
				_, err := fmt.Sscanf(yearVal, "%d", &year)
				if err != nil {
					http.Error(w, fmt.Sprintf("Invalid year: %v", err), http.StatusBadRequest)
					return
				}
			}

			// Execute build process
			builtData, _, err := build.Run(records, cache, p, opts)
			if err != nil {
				http.Error(w, fmt.Sprintf("Build run failed: %v", err), http.StatusInternalServerError)
				return
			}

			stats := recap.ComputeStats(builtData, year)

			// Response
			resp := map[string]interface{}{
				"recap": stats,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				log.Printf("Failed to encode response: %v", err)
			}
		})

		addr := fmt.Sprintf(":%d", servePort)
		log.Printf("Listening on %s...", addr)
		return http.ListenAndServe(addr, nil)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "server port")
	serveCmd.Flags().StringVar(&serveCacheDir, "cache-dir", store.DefaultCacheDir(), "metadata cache directory")
}
