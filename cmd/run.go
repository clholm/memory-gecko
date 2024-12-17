package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clholm/memory-gecko/server"
	"github.com/clholm/memory-gecko/youtube"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "runs both gather and serve functionality to display videos",
	Long:  `starts the web server and searches for videos to display at {host}:{port}/index.html`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// get api key
		if apiKey == "" {
			apiKey = viper.GetString("api-key")
		}
		if apiKey == "" {
			return fmt.Errorf("api-key is required either via --api-key flag or in config file")
		}

		// get host/port
		if host == "" {
			host = viper.GetString("host")
		}
		if port == "" {
			port = viper.GetString("port")
		}
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "8080"
		}

		// create context with cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// create channel for search results
		resultsChan := make(chan []youtube.SearchResult, 1)

		// start search in goroutine
		go func() {
			results, err := youtube.SearchVideos(ctx, apiKey)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error searching videos: %v\n", err)
				close(resultsChan)
				return
			}
			fmt.Printf("youtube search complete, found %d videos\n", len(results))
			if len(results) == 0 {
				fmt.Println("warning: no videos found in search")
			}
			resultsChan <- results
		}()

		// wait for results and start server
		results, ok := <-resultsChan
		if !ok {
			return fmt.Errorf("failed to get search results")
		}

		// start server with search results
		err := server.Run(
			ctx,
			os.Stdin,
			os.Stdout,
			os.Stderr,
			host,
			port,
			results,
		)

		if err != nil {
			return fmt.Errorf("server error: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// add flags from both gather and serve commands
	runCmd.Flags().StringVar(&apiKey, "api-key", "", "youtube api key (can also be set in config file)")
	runCmd.Flags().StringVarP(&host, "host", "H", "", "host to run the server (default: localhost)")
	runCmd.Flags().StringVarP(&port, "port", "p", "", "port to listen on (default: 8080)")

	// bind flags to viper
	viper.BindPFlag("api-key", runCmd.Flags().Lookup("api-key"))
	viper.BindPFlag("host", runCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", runCmd.Flags().Lookup("port"))
}
