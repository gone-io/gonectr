package create

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/gone-io/gonectr/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 克隆模板项目到指定路径
func cloneTemplate(templateURL, targetPath string) error {
	_, err := git.PlainClone(targetPath, false, &git.CloneOptions{
		URL:      templateURL,
		Progress: os.Stdout,
		Depth:    1,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	gitDir := filepath.Join(targetPath, ".git")
	err = os.RemoveAll(gitDir)
	if err != nil {
		return fmt.Errorf("failed to remove .git directory: %w", err)
	}

	return nil
}

var templateName, moduleName string

var supportedTemplates = []string{"web", "web+mysql"}
var supportedTemplatesMap map[string]bool

func init() {
	supportedTemplatesMap = make(map[string]bool)
	for _, template := range supportedTemplates {
		supportedTemplatesMap[template] = true
	}

	// 添加命令行标志
	Command.Flags().StringVarP(
		&templateName,
		"template-name",
		"t",
		"web",
		fmt.Sprintf("support template names: %s", strings.Join(supportedTemplates, ", ")),
	)

	Command.Flags().StringVarP(
		&moduleName,
		"module-name",
		"m",
		"",
		"module name",
	)
}

var Command = &cobra.Command{
	Use:   "create",
	Short: "Create a new Gone Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !supportedTemplatesMap[templateName] {
			return errors.New("unsupported template name")
		}

		if len(args) != 1 {
			return errors.New("please input project name or project path")
		}

		project := args[0]

		if moduleName == "" {
			moduleName = project
		}
		templateBaseUrl := "https://github.com/gone-io/template"
		if utils.IsInChina() {
			templateBaseUrl = "https://gitee.com/gone-io/template"
		}

		templateName := strings.Replace(templateName, "+", "-", -1)

		err := cloneTemplate(fmt.Sprintf("%s-%s", templateBaseUrl, templateName), project)
		if err != nil {
			return err
		}

		err = replaceModuleName(project, "template_module", moduleName)
		if err != nil {
			return err
		}
		err = os.Chdir(project)
		if err != nil {
			return err
		}

		command := exec.Command("go",
			[]string{
				"mod",
				"tidy",
			}...,
		)

		output, err := command.CombinedOutput()
		if err != nil {
			return err
		}
		rst := string(output)
		if rst != "" {
			fmt.Println(rst)
		}

		return nil
	},
}

// replaceModuleName 替换 go.mod 文件和所有 Go 源文件中的模块名称
func replaceModuleName(rootDir, oldModule, newModule string) error {
	// 1. 替换 go.mod 中的模块名
	modFile := filepath.Join(rootDir, "go.mod")
	err := replaceInFile(modFile, oldModule, newModule)
	if err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// 2. 遍历 Go 源文件并替换导入路径
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 Go 源文件
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

// replaceInFile 替换文件中的内容
func replaceInFile(filePath, oldModule, newModule string) error {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// 替换模块名称
	replacedData := bytes.ReplaceAll(data, []byte(oldModule), []byte(newModule))

	// 如果没有修改，返回
	if bytes.Equal(data, replacedData) {
		return nil
	}

	// 将修改后的数据写回文件
	err = ioutil.WriteFile(filePath, replacedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	fmt.Printf("Updated file: %s\n", filePath)
	return nil
}
