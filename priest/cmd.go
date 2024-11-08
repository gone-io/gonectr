package priest

import (
	"github.com/spf13/cobra"
)

var scanDirs []string
var packageName string
var functionName string
var outputFilePath string
var isStat bool
var isWatch bool

var Command = &cobra.Command{
	Use:   "priest",
	Short: "generate priest function",
	Long:  "gonectr priest -s ${scanPackageDir} -p ${pkgName} -f ${funcName} -o ${outputFilePath} [-w]",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return doAction(
			scanDirs,
			packageName,
			functionName,
			outputFilePath,
			isStat,
			isWatch,
		)
	},
}

func init() {
	Command.Flags().StringSliceVarP(&scanDirs, "scan-dir", "s", nil, "scan package dir")
	Command.Flags().StringVarP(&packageName, "package", "p", "", "package name of generated code")
	Command.Flags().StringVarP(&functionName, "function", "f", "", "function name of generated code")
	Command.Flags().StringVarP(&outputFilePath, "output", "o", "", "output filepath of generated code")
	Command.Flags().BoolVarP(&isStat, "stat", "t", false, "is stat process time")
	Command.Flags().BoolVarP(&isWatch, "watch", "w", false, "watch files change, and generate code when any files changed")
}
