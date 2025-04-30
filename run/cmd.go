package run

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gone-io/gonectl/utils"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "run",
	Short: "run gone project",
	Long: "This command will call `gonectl generate ...` to generate gone helper code first, and call `go run` to run gone project. " +
		"You can run `go help run` for looking up arguments.",
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

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if generatePath != "" {
		fmt.Printf("Find gonectl generate in `%s:%d`\n`%s`\n", generatePath, generateNumber, generateCommand)
		fmt.Printf("-> Change Dir to: `%s`\n", info.ModulePath)
		err = os.Chdir(info.ModulePath)
		if err != nil {
			return err
		}

		fmt.Printf("-> Execute `go generate ./...`\n")
		command := exec.Command("go", "generate", "./...")
		output, err := command.CombinedOutput()
		if err != nil {
			return err
		}
		println(string(output))

		fmt.Printf("-> Change Dir to: `%s`\n", workDir)
		err = os.Chdir(workDir)
		if err != nil {
			return err
		}
	} else {
		mainDir, err := filepath.Abs(packageName)
		if err != nil {
			return err
		}
		if strings.HasSuffix(mainDir, ".go") {
			mainDir = path.Dir(mainDir)
		}

		fmt.Printf("-> Execute `generate %s %s`\n", fmt.Sprintf("-s=%s", info.ModulePath), fmt.Sprintf("-m=%s", mainDir))
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
		println(string(output))
	}

	return utils.Command("go", append(
		[]string{goSubcommand},
		args...,
	))
}
