package create

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var templateName,
	moduleName string

var cacheDir string
var isListExample bool

func init() {
	Command.Flags().StringVarP(
		&cacheDir,
		"cache-dir",
		"c",
		"~/.gonectl",
		"cache dir",
	)

	// Add command line flags
	Command.Flags().StringVarP(
		&templateName,
		"template",
		"t",
		"simple",
		"template name; git repo url, or goner examples project, like: `mcp/stdio`, which can be listed by `gonectl create -l`",
	)

	Command.Flags().StringVarP(
		&moduleName,
		"module",
		"m",
		"",
		"module name for new project",
	)
	Command.Flags().BoolVarP(
		&isListExample,
		"ls",
		"l",
		false,
		"list all examples projects",
	)
}

var Command = &cobra.Command{
	Use:   "create <project dir>",
	Short: "Create a new Gone Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		cacheDir = strings.ReplaceAll(cacheDir, "~", usr.HomeDir)

		if isListExample {
			return listExamples()
		}

		if len(args) < 1 {
			return errors.New("please input project path")
		}

		return createProjectFromTpl(templateName, moduleName, args[0])
	},
}

// replaceModuleName replaces the module name in go.mod file and all Go source files
func replaceModuleName(rootDir, oldModule, newModule string) error {
	// 1. Replace module name in go.mod
	modFile := filepath.Join(rootDir, "go.mod")
	err := replaceInFile(modFile, oldModule, newModule)
	if err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// 2. Traverse Go source files and replace import paths
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process Go source files and XML files
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".xml") {
			err := replaceInFile(path, oldModule, newModule)
			if err != nil {
				return fmt.Errorf("failed to update file %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to update Go source files: %w", err)
	}

	return nil
}

// replaceInFile replaces content in the file
func replaceInFile(filePath, oldModule, newModule string) error {
	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Replace module name
	replacedData := bytes.ReplaceAll(data, []byte(oldModule), []byte(newModule))

	// Return if no changes were made
	if bytes.Equal(data, replacedData) {
		return nil
	}

	// Write the modified data back to file
	err = os.WriteFile(filePath, replacedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	fmt.Printf("Updated file: %s\n", filePath)
	return nil
}
