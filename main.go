package main

import (
	"fmt"
	"os"

	"github.com/gone-io/gonectr/generate"
	"github.com/gone-io/gonectr/mock"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gonectr",
	Short: "gonectr are tools for gone",
	Long:  `gonectr instructions`,
	Run: func(cmd *cobra.Command, args []string) {
		println("gonectr")
	},
}

func init() {
	rootCmd.AddCommand(mock.Command)
	rootCmd.AddCommand(generate.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
