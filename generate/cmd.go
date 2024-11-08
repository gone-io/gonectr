package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gone-io/gonectr/utils"

	"github.com/spf13/cobra"
)

type ModuleInfo struct {
	ModuleName string
	ModulePath string
}

var scanDirs []string
var mainPackageDir string
var preparerCode string
var preparerPackage string
var mainPackageName string
var excludeGoner []string

var Command = &cobra.Command{
	Use:   "generate",
	Short: "generate gone loading code and import code",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(scanDirs) == 0 {
			return errors.New("scan dir is empty")
		}

		if mainPackageDir == "" {
			mainPackageDir, err = findFirstMainPackageDir(scanDirs)
			if err != nil {
				return err
			}
		}

		moduleInfo, err := utils.FindModuleInfo(mainPackageDir)
		if err != nil {
			return err
		}

		var needImportPackages []string

		for _, dir := range scanDirs {
			packages, err := scanDirGenCode(dir, moduleInfo)
			if err != nil {
				return err
			}

			needImportPackages = append(needImportPackages, packages...)
		}

		if len(needImportPackages) > 0 {

			return genImportCode(mainPackageDir, needImportPackages)
		}
		return nil

	},
}

func init() {
	Command.Flags().StringSliceVarP(
		&scanDirs,
		"scan-dir",
		"s",
		[]string{"."},
		"scan dirs",
	)

	Command.Flags().StringVarP(
		&mainPackageDir,
		"main_package_dir",
		"m",
		"",
		"main package dir",
	)

	Command.Flags().StringVarP(
		&preparerCode,
		"preparer-code",
		"p",
		"gone.Default",
		"preparer code",
	)

	Command.Flags().StringVarP(
		&preparerPackage,
		"preparer-package",
		"r",
		"github.com/gone-io/gone",
		"preparer package",
	)

	Command.Flags().StringVarP(
		&mainPackageName,
		"main-package-name",
		"a",
		"main",
		"main package name",
	)

	Command.Flags().StringSliceVarP(&excludeGoner,
		"exclude-goner",
		"e",
		nil,
		"exclude goner",
	)
}

func isExclude(goneName string) bool {
	for _, name := range excludeGoner {
		var re = regexp.MustCompile(name)
		if re.MatchString(goneName) {
			return true
		}
	}
	return false
}

func findFirstMainPackageDir(scanDirs []string) (string, error) {
	var mainPackagePath string
	for _, dir := range scanDirs {
		// 遍历目录，查找 Go 文件
		err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 只处理 Go 文件
			if !info.IsDir() && filepath.Ext(file) == ".go" {
				// 检查是否为 main 包
				if isMainPackage(file) {
					mainPackagePath = path.Dir(file) // 记录找到的文件路径
					return filepath.SkipDir          // 停止遍历
				}
			}
			return nil
		})
		if err != nil && !errors.Is(err, filepath.SkipDir) {
			return "", err
		}
		if mainPackagePath != "" {
			return mainPackagePath, nil
		}

	}

	return "", fmt.Errorf("no main package found")
}

// isMainPackage 解析文件并检查是否包含 main 包
func isMainPackage(filePath string) bool {
	name, err := getModuleName(filePath)
	if err != nil {
		return false
	}
	return name == "main"
}

func getModuleName(filePath string) (string, error) {
	node, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		log.Println("Error parsing file:", filePath, err)
		return "", nil
	}
	return node.Name.Name, nil
}

func scanDirGenCode(dir string, moduleInfo *utils.ModuleInfo) ([]string, error) {
	relPath, err := filepath.Rel(moduleInfo.ModulePath, dir)
	if err != nil {
		return nil, err
	}
	if relPath == "." {
		relPath = ""
	} else if strings.HasPrefix(relPath, "..") {
		return nil, fmt.Errorf("scan dir %s is not in module %s", dir, moduleInfo.ModulePath)
	}

	packagePath := []string{moduleInfo.ModuleName}
	if relPath != "" {
		packagePath = append(packagePath,
			strings.Split(relPath, string(os.PathSeparator))...,
		)
	}

	return genLoadCodeForPackage(packagePath, dir)
}

func genLoadCodeForPackage(currentPackagePath []string, dir string) ([]string, error) {
	list, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var goFiles []string
	var subDirs []string
	for _, f := range list {
		if f.IsDir() {
			subDirs = append(subDirs, path.Join(dir, f.Name()))
		} else if strings.HasSuffix(f.Name(), ".go") &&
			!strings.HasSuffix(f.Name(), "_test.go") &&
			!strings.HasSuffix(f.Name(), ".gone.go") {
			goFiles = append(goFiles, path.Join(dir, f.Name()))
		}
	}
	var goners []string
	var tmpPackageName string

	if len(goFiles) > 0 {
		for _, filename := range goFiles {
			packageName, gonerModules, err := scanGoFile(filename)
			if err != nil {
				return nil, err
			}

			if tmpPackageName == "" {
				tmpPackageName = packageName
			}

			if tmpPackageName != packageName {
				return nil, fmt.Errorf("package name %s is not equal to %s", packageName, tmpPackageName)
			}

			goners = append(goners, gonerModules...)
		}
	}

	gonersList := goners
	goners = nil
	for _, g := range gonersList {
		if !isExclude(g) {
			goners = append(goners, g)
		}
	}

	var needImportPackages []string
	if len(goners) > 0 {
		if len(currentPackagePath) > 1 && currentPackagePath[len(currentPackagePath)-1] != tmpPackageName {
			currentPackagePath = append(currentPackagePath[0:len(currentPackagePath)-1], tmpPackageName)
		}

		err = genLoadCode(goners, currentPackagePath[len(currentPackagePath)-1], dir)
		if err != nil {
			return nil, err
		}

		needImportPackages = append(needImportPackages, strings.Join(currentPackagePath, "/"))
	}

	for _, subDir := range subDirs {
		base := path.Base(subDir)
		packages, err := genLoadCodeForPackage(append(currentPackagePath, base), subDir)
		if err != nil {
			return nil, err
		}
		needImportPackages = append(needImportPackages, packages...)
	}
	return needImportPackages, nil
}

const loadTpl = utils.GenerateBy + `
package %s

import "%s"

func init() {
	%s%s
}
`

func genLoadCode(goners []string, packageName string, packageDir string) error {
	loadCode := ""
	for _, goner := range goners {
		loadCode = fmt.Sprintf("%s.\n\t\tLoad(&%s{})", loadCode, goner)
	}

	code := fmt.Sprintf(loadTpl, packageName, preparerPackage, preparerCode, loadCode)
	return os.WriteFile(path.Join(packageDir, "init.gone.go"), []byte(code), 0644)
}

func scanGoFile(filename string) (string, []string, error) {
	// 打开文件
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	// 获取包名
	packageName := file.Name.Name

	// 用于存储符合条件的结构体名称
	var structNames []string

	// 遍历 AST 节点
	for _, decl := range file.Decls {
		// 检查是否为通用声明（GenDecl）
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// 遍历类型声明
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 检查类型声明是否为结构体
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 检查结构体字段是否嵌入了 gone.Flag
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					// 匿名字段，可能是嵌入
					ident, ok := field.Type.(*ast.SelectorExpr)
					if ok && ident.Sel.Name == "Flag" {
						if x, ok := ident.X.(*ast.Ident); ok && x.Name == "gone" {
							structNames = append(structNames, typeSpec.Name.Name)
							break
						}
					}
				}
			}
		}
	}

	return packageName, structNames, nil
}

const importTPl = utils.GenerateBy + `

package %s

import (
%s
)`

func genImportCode(mainPackageDir string, needImportPackages []string) error {

	var imports []string
	for _, pkg := range needImportPackages {
		imports = append(imports, fmt.Sprintf("\t_ \"%s\"", pkg))
	}
	code := fmt.Sprintf(importTPl, mainPackageName, strings.Join(imports, "\n"))

	return os.WriteFile(path.Join(mainPackageDir, "import.gone.go"), []byte(code), 0644)
}

//func findModuleInfo(dir string) (*ModuleInfo, error) {
//	modulePath, err := findGoModFile(dir)
//	if err != nil {
//		return nil, err
//	}
//
//	moduleName, err := parseModuleName(path.Join(modulePath, "go.mod"))
//	if err != nil {
//		return nil, err
//	}
//
//	return &ModuleInfo{
//		ModuleName: moduleName,
//		ModulePath: modulePath,
//	}, nil
//}

//// findGoModFile 从指定目录向上逐层搜索 "go.mod" 文件
//func findGoModFile(dir string) (string, error) {
//	for {
//		goModPath := filepath.Join(dir, "go.mod")
//
//		// 检查当前目录是否有 "go.mod" 文件
//		if _, err := os.Stat(goModPath); err == nil {
//			return filepath.Dir(goModPath), nil
//		}
//
//		// 获取上级目录
//		parentDir := filepath.Dir(dir)
//
//		// 如果已经到达根目录，就退出
//		if parentDir == dir {
//			return "", fmt.Errorf("未找到 go.mod 文件")
//		}
//
//		// 更新目录为上级目录，继续搜索
//		dir = parentDir
//	}
//}
//
//// parseModuleName 读取 go.mod 文件并解析出 module 名称
//func parseModuleName(goModPath string) (string, error) {
//	file, err := os.Open(goModPath)
//	if err != nil {
//		return "", fmt.Errorf("无法打开文件: %w", err)
//	}
//	defer func(file *os.File) {
//		err := file.Close()
//		if err != nil {
//			fmt.Println("关闭文件出错:", err)
//		}
//	}(file)
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line := strings.TrimSpace(scanner.Text())
//
//		// 跳过空行和注释行
//		if line == "" || strings.HasPrefix(line, "//") {
//			continue
//		}
//
//		// 检查行是否以 "module" 开头
//		if strings.HasPrefix(line, "module ") {
//			// 提取模块名称
//			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module "))
//			return moduleName, nil
//		}
//	}
//
//	if err := scanner.Err(); err != nil {
//		return "", fmt.Errorf("读取文件出错: %w", err)
//	}
//
//	return "", fmt.Errorf("未找到 module 声明")
//}
