package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/kmdkuk/nfrecap/internal/recap"
)

var (
	recapIn   string
	recapOut  string
	recapYear int
)

var recapCmd = &cobra.Command{
	Use:   "recap",
	Short: "Generate stats-heavy recap markdown from built JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		year := recapYear
		if year == 0 {
			year = time.Now().Year()
		}

		built, err := recap.ReadBuiltJSON(recapIn)
		if err != nil {
			return err
		}

		stats := recap.ComputeStats(built, year)
		md := recap.RenderMarkdown(stats)

		if recapOut == "-" {
			fmt.Print(md)
			return nil
		}
		return os.WriteFile(recapOut, []byte(md), 0644)
	},
}

func init() {
	rootCmd.AddCommand(recapCmd)

	recapCmd.Flags().StringVarP(&recapIn, "in", "i", "", "input built JSON file (from `nfrecap build`)")
	recapCmd.Flags().StringVarP(&recapOut, "out", "o", "-", "output markdown file ('-' for stdout)")
	recapCmd.Flags().IntVarP(&recapYear, "year", "y", 0, "target year (default: current year)")

	_ = recapCmd.MarkFlagRequired("in")
}
