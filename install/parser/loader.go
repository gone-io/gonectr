package parser

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gone-io/gonectr/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

func New(workDir string, module string) (*LoaderParser, error) {
	loaderParser := LoaderParser{
		workDir: workDir,
		module:  module,
	}
	return &loaderParser, loaderParser.Init()
}

type Import struct {
	Alias   string
	PkgName string
	PkgID   string
}

type LoadFunc struct {
	Name    string
	PkgID   string
	PkgName string
}

func (l LoadFunc) ID() string {
	return fmt.Sprintf("%s.%s", l.PkgID, l.Name)
}

func (l LoadFunc) String() string {
	return fmt.Sprintf("%s.%s", l.PkgID, l.Name)
}

type LoaderParser struct {
	workDir string
	module  string

	currentModule        string
	currentModuleAbsPath string
	loaderFile           string

	fset      *token.FileSet
	node      *ast.File
	Imports   []*Import
	loadFuncs map[string]*LoadFunc
}

const GoneModule = "github.com/gone-io/gone/v2"
const LoaderFile = "module.load.go"

func (s *LoaderParser) Init() error {
	defer utils.TimeStat("Init")()

	s.loadFuncs = make(map[string]*LoadFunc)
	info, err := utils.FindModuleInfo(s.workDir)
	if err != nil {
		return err
	}
	s.currentModule = info.ModuleName
	s.currentModuleAbsPath = info.ModulePath
	s.loaderFile = path.Join(s.currentModuleAbsPath, LoaderFile)
	return s.checkOrGenerate()
}

func (s *LoaderParser) checkOrGenerate() error {
	fileExists := false
	fileInfo, err := os.Stat(s.loaderFile)
	if err == nil {
		if fileInfo.IsDir() {
			return fmt.Errorf("%s is a dir: %w", s.loaderFile, err)
		}
		fileExists = true
	}
	if !fileExists {
		packageName := GetDirPackageName(s.currentModuleAbsPath)
		if packageName == "" {
			packageName = path.Base(s.currentModule)
		}

		content := fmt.Sprintf(codeTpl,
			utils.GenerateBy,
			packageName,
		)
		if err = os.WriteFile(s.loaderFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("create %s failed: %v", s.loaderFile, err)
		}
	}
	return s.ast()
}
func (s *LoaderParser) ast() (err error) {
	defer utils.TimeStat("Create AST Tree for loader.gone.go")()
	s.fset = token.NewFileSet()
	s.node, err = parser.ParseFile(s.fset, s.loaderFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse %s failed: %v", s.loaderFile, err)
	}
	return nil
}

func (s *LoaderParser) ParseModuleLoader() (loaders []*LoadFunc, err error) {
	defer utils.TimeStat("ParseModuleLoader")()

	var absPaths []string
	for _, m := range []string{GoneModule, s.module} {
		absPath, err := GetDependentModuleAbsPath(m)
		if err != nil {
			return nil, err
		}
		absPaths = append(absPaths, absPath)
	}

	fmt.Printf("\tloaders module info ...\n")
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}, absPaths...)

	if err != nil {
		return nil, err
	}

	fmt.Printf("\tlooking for LoadFunc ...\n")

	var loaderType types.Type
	for _, pkg := range pkgs {
		if pkg.ID == GoneModule {
			if obj := pkg.Types.Scope().Lookup("Loader"); obj != nil {
				loaderType = obj.Type()
				break
			}
		}
	}

	for _, pkg := range pkgs {
		if pkg.ID == GoneModule {
			continue
		}

		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				fnDecl, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}

				obj := pkg.TypesInfo.Defs[fnDecl.Name]
				if obj == nil {
					return true
				}
				sig, ok := obj.Type().(*types.Signature)
				if !ok {
					return true
				}

				if sig.Params().Len() == 1 &&
					sig.Results().Len() == 1 &&
					types.Identical(sig.Params().At(0).Type(), loaderType) &&
					sig.Results().At(0).Type().String() == "error" {

					loaders = append(loaders, &LoadFunc{
						PkgID:   pkg.ID,
						Name:    fnDecl.Name.Name,
						PkgName: path.Base(pkg.ID),
					})
				}
				return true
			})
		}
	}
	return
}

func (s *LoaderParser) ParseImports() {
	defer utils.TimeStat("ParseImports")()

	for _, imp := range s.node.Imports {
		if imp.Path != nil {
			var name string

			if imp.Name != nil {
				name = imp.Name.Name
			}

			s.Imports = append(s.Imports, &Import{
				Alias: name,
				PkgID: strings.Trim(imp.Path.Value, `"`),
			})
		}
	}
}

func (s *LoaderParser) ParseLoadFuncs() {
	defer utils.TimeStat("ParseLoadFuncs")()

	m := convertor.ToMap(s.Imports, func(t *Import) (string, *Import) {
		if t.Alias == "" {
			return path.Base(t.PkgID), t
		}
		return t.Alias, t
	})

	for _, decl := range s.node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok && len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == "loaders" {
					if len(valueSpec.Values) > 0 {
						if arrayLit, ok := valueSpec.Values[0].(*ast.CompositeLit); ok {
							for _, elt := range arrayLit.Elts {
								switch expr := elt.(type) {
								case *ast.SelectorExpr:
									if x, ok := expr.X.(*ast.Ident); ok {
										if imp, ok := m[x.Name]; ok {
											loadFunc := LoadFunc{
												PkgID:   imp.PkgID,
												PkgName: path.Base(imp.PkgID),
												Name:    expr.Sel.Name,
											}
											s.loadFuncs[loadFunc.ID()] = &loadFunc
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (s *LoaderParser) Select(options []*LoadFunc, cmdSelected []string) error {
	defer utils.TimeStat("Select")()

	if len(cmdSelected) > 0 {
		for _, selected := range cmdSelected {

			find := false
			for _, option := range options {
				if selected == option.Name ||
					selected == fmt.Sprintf("%s.%s", option.PkgName, option.Name) ||
					selected == fmt.Sprintf("%s.%s", option.PkgID, option.Name) {
					s.loadFuncs[option.ID()] = option
					find = true
					break
				}
			}
			if !find {
				var tmp []string
				for _, o := range options {
					tmp = append(tmp, fmt.Sprintf("\t - %s", o.String()))
				}
				return fmt.Errorf("cannot select %s, only find:\n %s\n", selected, strings.Join(tmp, "\n"))
			}
		}
		return nil
	}
	return s.userSelect(options)
}

func (s *LoaderParser) userSelect(options []*LoadFunc) error {
	var defaultVal []string
	type Select struct {
		loadFunc *LoadFunc
		selected bool
	}

	selectMap := convertor.ToMap(options, func(option *LoadFunc) (string, *Select) {
		_, ok := s.loadFuncs[option.ID()]
		if ok {
			defaultVal = append(defaultVal, option.ID())
		}
		return option.ID(), &Select{
			loadFunc: option,
			selected: false,
		}
	})

	if len(options) == 1 && len(defaultVal) != 1 {
		s.loadFuncs[options[0].ID()] = options[0]
		return nil
	}

	prompt := &survey.MultiSelect{
		Message: "Add or Remove goner LoadFunc:",
		Options: slice.Map(options, func(i int, v *LoadFunc) string {
			return v.ID()
		}),
		Default:  defaultVal,
		PageSize: 5,
	}
	var selected []string
	if err := survey.AskOne(prompt, &selected); err != nil {
		return fmt.Errorf("select error:%v", err)
	}

	for _, id := range selected {
		selectMap[id].selected = true
	}
	for id, selectItem := range selectMap {
		if selectItem.selected {
			s.loadFuncs[id] = selectMap[id].loadFunc
		} else {
			delete(s.loadFuncs, id)
		}
	}
	return nil
}

func (s *LoaderParser) GenerateCode() (importCode, loadFuncCode string) {
	defer utils.TimeStat("Generate Code")()
	m := make(map[string]*Import)
	pkgNameMap := make(map[string]*Import)

	for _, l := range s.loadFuncs {
		if _, ok := m[l.PkgID]; !ok {
			pkgName := path.Base(l.PkgID)
			i := Import{
				PkgID:   l.PkgID,
				PkgName: pkgName,
			}

			var alias string
			if _, ok := pkgNameMap[pkgName]; ok {
				alias = generateNotDuplicateAlias(pkgNameMap, pkgName)
				pkgNameMap[alias] = &i
				i.Alias = alias
			} else {
				pkgNameMap[pkgName] = &i
			}
			m[l.PkgID] = &i
		}
	}
	pkgNames := make([]string, 0, len(pkgNameMap))
	for _, i := range pkgNameMap {
		pkgNames = append(pkgNames, i.PkgName)
	}
	sort.Strings(pkgNames)

	importCode = "\t\"github.com/gone-io/gone/v2\"\n\t\"github.com/gone-io/goner/g\""
	for _, pkgName := range pkgNames {
		i := pkgNameMap[pkgName]
		if i.Alias != "" {
			importCode += fmt.Sprintf("\n\t%s \"%s\"", i.Alias, i.PkgID)
		} else {
			importCode += fmt.Sprintf("\n\t\"%s\"", i.PkgID)
		}
	}

	var LoadFuncStr []string

	for _, l := range s.loadFuncs {
		i := m[l.PkgID]
		if i.Alias != "" {
			LoadFuncStr = append(LoadFuncStr, fmt.Sprintf("%s.%s", i.Alias, l.Name))
		} else {
			LoadFuncStr = append(LoadFuncStr, fmt.Sprintf("%s.%s", i.PkgName, l.Name))
		}
	}
	sort.Strings(LoadFuncStr)
	for _, str := range LoadFuncStr {
		loadFuncCode += fmt.Sprintf("\n\t%s,", str)
	}

	importCode = fmt.Sprintf("import(\n%s\n)", importCode)
	loadFuncCode = fmt.Sprintf("var loaders = []gone.LoadFunc{%s\n}", loadFuncCode)
	return
}

func (s *LoaderParser) read() (string, error) {
	content, err := os.ReadFile(s.loaderFile)
	if err != nil {
		return "", fmt.Errorf("read file %s failed: %v", s.loaderFile, err)
	}
	return string(content), nil
}

func (s *LoaderParser) Save(importCode, loadFuncCode string) error {
	defer utils.TimeStat("Save Code")()
	if fileContent, err := s.read(); err != nil {
		return err
	} else {
		fileContent = s.replaceCode(fileContent, importCode, loadFuncCode)
		return os.WriteFile(s.loaderFile, []byte(fileContent), 0644)
	}
}

func (s *LoaderParser) replaceCode(fileContent, importCode, loadFuncCode string) string {
	// 使用正则表达式替换import部分
	importRegex := regexp.MustCompile(`(?s)import\s*\(([^\)]*)\)`)
	fileContent = importRegex.ReplaceAllString(fileContent, importCode)

	// 使用正则表达式替换loaders变量定义部分
	loadersRegex := regexp.MustCompile(`(?s)var\s+loaders\s*=\s*\[\]gone\.LoadFunc\s*\{[^\}]*\}`)
	fileContent = loadersRegex.ReplaceAllString(fileContent, loadFuncCode)

	return fileContent
}

func (s *LoaderParser) Execute(cmdSelected []string, onlyPrint bool) error {
	s.ParseImports()
	s.ParseLoadFuncs()
	loaders, err := s.ParseModuleLoader()
	if err != nil {
		return err
	}
	if onlyPrint {
		fmt.Printf("loaders in %s\n", s.module)
		for _, l := range loaders {
			fmt.Printf("- %s", l.String())
		}
		return nil
	}
	if err = s.Select(loaders, cmdSelected); err != nil {
		return err
	}

	importCode, loadFuncCode := s.GenerateCode()
	return s.Save(importCode, loadFuncCode)
}
