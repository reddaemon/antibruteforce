package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	subnet  string
	address string
	login   string
	ip      string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "antibruteforce",
	Short: "antibruteforce",
	Long:  `antibruteforce`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Use antibruteforce [command]\nRun 'antibruteforce --help' for usage.\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("unable to execute: %v", err)
	}
}
