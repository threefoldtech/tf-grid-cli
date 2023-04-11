// Package cmd for parsing command line arguments
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tf-grid",
	Short: "A cli for interacting with Threefold Grid",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(getCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}
