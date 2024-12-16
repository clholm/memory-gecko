package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "memory-gecko",
	Short: "memory-gecko helps you watch YouTube videos that follow iPhone's default naming convention.",
	Long: `memory-gecko helps you watch YouTube videos that follow iPhone's default naming convention, 
	IMG_XXXX. Inspired by https://ben-mini.github.io/2024/img-0416,
	memory-gecko runs as a web application. 
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// execute adds all child commands to the root command and sets flags appropriately.
// this is called by main.main(). it only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.memory-gecko)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// search config in home directory with name ".memory-gecko" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".memory-gecko")
	}

	viper.AutomaticEnv() // read in environment variables that match flags

	// if a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
