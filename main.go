package main

import (
	"fmt"
	"github.com/gone-io/gonectr/create"
	"os"

	"github.com/gone-io/gonectr/build"
	"github.com/gone-io/gonectr/priest"
	"github.com/gone-io/gonectr/run"

	"github.com/gone-io/gonectr/generate"
	"github.com/gone-io/gonectr/mock"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "gonectr",
	Short: "gonectr is a tool for gone project",
	Long: `gonectr is a command-line tool designed for generating Gone projects
and serving as a utility for Gone projects, such as code generation,
compilation, and running Gone projects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			fmt.Printf("Gonectr version: %s\n", version)
			return nil
		} else {
			return cmd.Help()
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "version", "v", false, "Show version")

	rootCmd.AddCommand(mock.Command)
	rootCmd.AddCommand(generate.Command)
	rootCmd.AddCommand(run.Command)
	rootCmd.AddCommand(build.Command)
	rootCmd.AddCommand(priest.Command)
	rootCmd.AddCommand(create.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
