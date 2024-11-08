package build

import (
	"github.com/gone-io/gonectr/run"
	"github.com/spf13/cobra"
	"os"
)

var Command = &cobra.Command{
	Use:   "build",
	Short: "build gone project",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return run.GenerateAndRunGoSubCommand("build", os.Args[2:])
	},
}

func init() {
	Command.FParseErrWhitelist.UnknownFlags = true
}
