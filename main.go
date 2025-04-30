package main

import (
	"fmt"
	"github.com/gone-io/gonectl/create"
	"github.com/gone-io/gonectl/install"
	"os"

	"github.com/gone-io/gonectl/build"
	"github.com/gone-io/gonectl/priest"
	"github.com/gone-io/gonectl/run"

	"github.com/gone-io/gonectl/generate"
	"github.com/gone-io/gonectl/mock"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "gonectl",
	Short: "gonectl is a tool for gone project",
	Long: `gonectl is a command-line tool designed for generating Gone projects
and serving as a utility for Gone projects, such as create project, install gone module, generate helper code,
compilation, and running gone project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			fmt.Printf("gonectl version: %s\n", version)
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
	rootCmd.AddCommand(install.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
