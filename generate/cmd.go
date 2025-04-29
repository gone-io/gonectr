package generate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
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
var goneVersion string

var Command = &cobra.Command{
	Use:   "generate",
	Short: "generate gone loading code and import code",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(scanDirs) == 0 {
			return errors.New("scan dir is empty")
		}

		//获取当前工作目录
		cwd, _ := os.Getwd()
		fmt.Printf("current work dir: %s\n", cwd)

		for i := range scanDirs {
			absolutePath, err := filepath.Abs(scanDirs[i])
			if err != nil {
				return err
			}
			scanDirs[i] = absolutePath
		}

		if mainPackageDir == "" {
			mainPackageDir, err = findFirstMainPackageDir(scanDirs)
			if err != nil {
				return err
			}
		}

		fmt.Printf("main package dir: %s\n", mainPackageDir)
		fmt.Printf("scan dirs: %v\n", scanDirs)

		moduleInfo, err := utils.FindModuleInfo(mainPackageDir)
		if err != nil {
			return err
		}

		var needImportPackages []string

		for _, dir := range scanDirs {
			packages, err := scanDirGenCode(dir, moduleInfo, mainPackageDir)
			if err != nil {
				return err
			}

			needImportPackages = append(needImportPackages, packages...)
		}

		_ = os.Remove(path.Join(mainPackageDir, ImportFileName))
		if len(needImportPackages) > 0 {
			return genImportCode(mainPackageDir, needImportPackages)
		}

		err = os.Chdir(moduleInfo.ModulePath)
		if err != nil {
			return err
		}

		return utils.Command("go", []string{"mod", "tidy"})
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
		"gone",
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

	Command.Flags().StringVarP(&goneVersion, "version", "v", "", "gone version")
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

func getGoneVersionFromModuleFile() string {
	if goneVersion == "" {
		goneVersion = utils.GetGoneVersionFromModuleFile(scanDirs, nil)
	}
	return goneVersion
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

func scanDirGenCode(dir string, moduleInfo *utils.ModuleInfo, mainPackageDir string) ([]string, error) {
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

	return genLoadCodeForPackage(packagePath, dir, mainPackageDir)
}

const ImportFileName = "import.gone.go"
const InitFileName = "init.gone.go"

func genLoadCodeForPackage(currentPackagePath []string, currentScanDir string, mainPackageDir string) ([]string, error) {
	_ = os.Remove(path.Join(currentScanDir, InitFileName))

	list, err := os.ReadDir(currentScanDir)
	if err != nil {
		return nil, err
	}

	var goFiles []string
	var subDirs []string
	for _, f := range list {
		if f.IsDir() {
			subDirs = append(subDirs, path.Join(currentScanDir, f.Name()))
		} else if strings.HasSuffix(f.Name(), ".go") &&
			!strings.HasSuffix(f.Name(), "_test.go") &&
			!strings.HasSuffix(f.Name(), ".gone.go") {
			goFiles = append(goFiles, path.Join(currentScanDir, f.Name()))
		}
	}
	var goners []string
	var loadFuncs []string
	var correctPackageName string

	if len(goFiles) > 0 {
		for _, filename := range goFiles {
			packageName, gonerStructNames, loadFuncNames, err := scanGoFile(filename, nil)
			if err != nil {
				return nil, err
			}

			if correctPackageName == "" {
				correctPackageName = packageName
			}

			if correctPackageName != packageName {
				return nil, fmt.Errorf("package name %s is not equal to %s", packageName, correctPackageName)
			}

			goners = append(goners, gonerStructNames...)
			loadFuncs = append(loadFuncs, loadFuncNames...)
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
	if len(goners) > 0 || len(loadFuncs) > 0 {
		filename, content := genLoadCode(goners, loadFuncs, correctPackageName, currentScanDir)
		if err = os.WriteFile(filename, []byte(content), 0644); err != nil {
			return nil, err
		}

		// 相同目录不生成到 import.gone.go
		if mainPackageDir != currentScanDir && correctPackageName != "main" {
			needImportPackages = append(needImportPackages, strings.Join(currentPackagePath, "/"))
		}
	}

	for _, subDir := range subDirs {
		base := path.Base(subDir)
		packages, err := genLoadCodeForPackage(append(currentPackagePath, base), subDir, mainPackageDir)
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

func genLoadCode(goners []string, loadFuncs []string, packageName string, packageDir string) (filename, content string) {
	loadCode := ""
	if len(loadFuncs) > 0 {
		for _, loadFunc := range loadFuncs {
			loadCode = fmt.Sprintf("%s.\n\t\tLoads(%s)", loadCode, loadFunc)
		}
	} else {
		for _, goner := range goners {
			loadCode = fmt.Sprintf("%s.\n\t\tLoad(&%s{})", loadCode, goner)
		}
	}

	getGoneVersionFromModuleFile()
	if preparerPackage == "github.com/gone-io/gone" && goneVersion != "v1" {
		preparerPackage = fmt.Sprintf("github.com/gone-io/gone/%s", goneVersion)
	}

	code := fmt.Sprintf(loadTpl, packageName, preparerPackage, preparerCode, loadCode)

	return path.Join(packageDir, InitFileName), code
}

func scanGoFile(filename string, src any) (packageName string, structNames, loadFuncNames []string, err error) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, src, parser.ParseComments)
	if err != nil {
		return "", nil, nil, err
	}

	// 获取包名
	packageName = file.Name.Name

	// 遍历 AST 节点
	for _, decl := range file.Decls {

		if funcDecl, ok := decl.(*ast.FuncDecl); ok &&
			funcDecl.Recv == nil &&
			funcDecl.Type.TypeParams == nil &&
			len(funcDecl.Type.Params.List) == 1 &&
			len(funcDecl.Type.Results.List) == 1 {

			paramName := types.ExprString(funcDecl.Type.Params.List[0].Type)
			returnName := types.ExprString(funcDecl.Type.Results.List[0].Type)
			if "gone.Loader" == paramName && "error" == returnName {
				loadFuncNames = append(loadFuncNames, funcDecl.Name.Name)
			}
			continue
		}

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

	return packageName, structNames, loadFuncNames, nil
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

	return os.WriteFile(path.Join(mainPackageDir, ImportFileName), []byte(code), 0644)
}
