package install

import (
	"errors"
	"github.com/gone-io/gonectr/install/parser"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var installModule string
var useLoaders []string
var onlyPrintLoadFunc bool

func init() {
	Command.Flags().BoolVarP(
		&onlyPrintLoadFunc,
		"test",
		"t",
		false,
		"only print `LoadFunc` list",
	)
}

var Command = &cobra.Command{
	Use:   "install <moduleName> [loadFuncName[,loadFuncName[,...]]]",
	Short: "install goner component and generate loaded code",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) < 1 {
			return errors.New("must provide package full name")
		}

		installModule = args[0]
		if len(args) > 1 {
			useLoaders = strings.Split(args[1], ",")
		}
		if err = Install(installModule, useLoaders, onlyPrintLoadFunc); err != nil {
			return
		}
		return nil
	},
}

func Install(moduleName string, loaderNames []string, onlyPrint bool) (err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	loaderParser, err := parser.New(workDir, moduleName)
	if err != nil {
		return err
	}
	return loaderParser.Execute(loaderNames, onlyPrint)
}
