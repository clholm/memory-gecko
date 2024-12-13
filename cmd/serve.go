package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/clholm/memory-gecko/server"
	"github.com/spf13/cobra"
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
	serveCmd.Flags().StringVarP(&host, "host", "H", "localhost", "host to bind the server")
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to run the server on")
}
