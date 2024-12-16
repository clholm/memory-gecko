package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clholm/memory-gecko/youtube"
)

var apiKey string

var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "search for videos matching the IMG_XXXX pattern",
	Long:  `searches youtube for videos that match iPhone's default video naming pattern (IMG_XXXX)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// if api-key wasn't provided via flag, try to get it from config
		if apiKey == "" {
			apiKey = viper.GetString("api-key")
		}

		// check if we have an API key from either source
		if apiKey == "" {
			return fmt.Errorf("api-key is required either via --api-key flag or in config file")
		}

		// create context
		ctx := context.Background()

		// perform search
		results, err := youtube.SearchVideos(ctx, apiKey)
		if err != nil {
			return fmt.Errorf("error searching videos: %w", err)
		}

		// write results as JSON to stdout
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(results); err != nil {
			return fmt.Errorf("error encoding results: %w", err)
		}

		return nil
	},
}

func init() {
	// add the gather command to the root command
	rootCmd.AddCommand(gatherCmd)

	// add api-key flag
	gatherCmd.Flags().StringVar(&apiKey, "api-key", "", "youtube api key (can also be set in config file)")

	// bind the flag to viper
	viper.BindPFlag("api-key", gatherCmd.Flags().Lookup("api-key"))
}
