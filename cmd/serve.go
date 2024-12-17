package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/clholm/memory-gecko/server"
)

var (
	host string
	port string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the memory-gecko server",
	Long:  `starts the memory-gecko web server with specified host and port`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// if host/port weren't provided via flags, try to get them from config
		if host == "" {
			host = viper.GetString("host")
		}
		if port == "" {
			port = viper.GetString("port")
		}

		// set defaults if still not configured
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "8080"
		}

		// get api key
		if apiKey == "" {
			apiKey = viper.GetString("api-key")
		}
		if apiKey == "" {
			return fmt.Errorf("api-key is required either via --api-key flag or in config file")
		}

		// initialize a context
		ctx := context.Background()

		// call the run function from the server package
		err := server.Run(
			ctx,
			os.Stdin,
			os.Stdout,
			os.Stderr,
			host,
			port,
			apiKey,
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	// add the serve command to the root command
	rootCmd.AddCommand(serveCmd)

	// define flags for host and port
	serveCmd.Flags().StringVarP(&host, "host", "H", "", "host to run the server (default: localhost)")
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "port to listen on (default: 8080)")

	// bind flags to viper
	viper.BindPFlag("host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}
