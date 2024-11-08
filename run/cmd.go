package run

import (
	"fmt"
	"github.com/gone-io/gonectr/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"strings"
)

var Command = &cobra.Command{
	Use:   "run",
	Short: "run gone project",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return GenerateAndRunGoSubCommand("run", os.Args[2:])
	},
}

func init() {
	Command.FParseErrWhitelist.UnknownFlags = true
}

func GenerateAndRunGoSubCommand(goSubcommand string, args []string) error {
	packageName := utils.ExtractPackageArg(args)
	info, err := utils.FindModuleInfo(packageName)
	if err != nil {
		return err
	}

	generatePath, generateNumber, generateCommand, err := utils.FindFirstGoGenerateLine(info.ModulePath)
	if err != nil {
		return err
	}

	if generatePath != "" {
		fmt.Printf("Find gonectr generate in `%s:%d`\n The line is `%s`\n", generatePath, generateNumber, generateCommand)
		thePath := fmt.Sprintf("%s/...", info.ModulePath)
		fmt.Printf("Execute `go generate %s`\n", thePath)

		command := exec.Command("go", "generate", thePath)
		output, err := command.CombinedOutput()
		if err != nil {
			return err
		}
		println(output)
	} else {
		mainDir := packageName
		if strings.HasSuffix(mainDir, ".go") {
			mainDir = path.Dir(mainDir)
		}

		fmt.Printf("execute `generate %s %s`", fmt.Sprintf("-s=%s", info.ModulePath), fmt.Sprintf("-m=%s", mainDir))
		command := exec.Command(
			os.Args[0],
			"generate",
			fmt.Sprintf("-s=%s", info.ModulePath),
			fmt.Sprintf("-m=%s", mainDir),
		)
		output, err := command.CombinedOutput()
		if err != nil {
			return err
		}
		println(output)
	}

	command := exec.Command(
		"go",
		append(
			[]string{goSubcommand},
			args...,
		)...,
	)
	_, err = command.CombinedOutput()

	return err
}
