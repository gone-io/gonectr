package parser

import (
	"fmt"
	"github.com/gone-io/gonectl/utils"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"strings"
)

func GetDependentModuleAbsPath(module string) (string, error) {
	fmt.Printf("\tgo get -u %s\n", module)
	err := utils.Command("go", []string{"get", "-u", module})
	if err != nil {
		return "", err
	}

	fmt.Printf("\tgo list -m -f \"{{.Dir}}\" %s\n", module)
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", module)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cannot get module abspath: %v, output: %s", err, string(output))
	}
	fmt.Printf("\tmodule(%s) abspath=> %s\n\n", module, strings.TrimSpace(string(output)))
	return strings.TrimSpace(string(output)), nil
}

// GetDirPackageName get package name of dir path
func GetDirPackageName(dir string) string {
	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") && !strings.HasSuffix(file.Name(), "_test.go") {
			fset := token.NewFileSet()
			filePath := path.Join(dir, file.Name())
			node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
			if err == nil && node != nil && node.Name != nil {
				return node.Name.Name
			}
		}
	}
	return ""
}

func generateNotDuplicateAlias(importMap map[string]*Import, name string) string {
	if _, ok := importMap[name]; !ok {
		return name
	}
	for i := 1; i < 1000; i++ {
		newName := fmt.Sprintf("%s%d", name, i)
		if _, ok := importMap[newName]; !ok {
			return newName
		}
	}
	panic("too many times to try generate name")
}
