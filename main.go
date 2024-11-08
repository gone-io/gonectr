package main

import (
	"fmt"
	"github.com/gone-io/gonectr/build"
	"github.com/gone-io/gonectr/run"
	"os"

	"github.com/gone-io/gonectr/generate"
	"github.com/gone-io/gonectr/mock"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gonectr",
	Short: "gonectr is a tool for gone project",
	Long:  `gonectr instructions`,
	Run: func(cmd *cobra.Command, args []string) {
		println("gonectr")
	},
}

func init() {
	rootCmd.AddCommand(mock.Command)
	rootCmd.AddCommand(generate.Command)
	rootCmd.AddCommand(run.Command)
	rootCmd.AddCommand(build.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
