package main

import (
	"fmt"
	"os"

	"github.com/gone-io/tools/cmd/goner/mock"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goner",
	Short: "goner is a tool for gone",
	Long:  `goner instructions`,
	Run: func(cmd *cobra.Command, args []string) {
		println("goner")
	},
}

func init() {
	rootCmd.AddCommand(mock.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
