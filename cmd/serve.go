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
	Run: func(cmd *cobra.Command, args []string) {
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
		)

		if err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// add the serve command to the root command
	rootCmd.AddCommand(serveCmd)

	// define flags for host and port
	serveCmd.Flags().StringVarP(&host, "host", "H", "", "host to bind the server (default: localhost)")
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "port to run the server on (default: 8080)")

	// bind flags to viper
	viper.BindPFlag("host", serveCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}
