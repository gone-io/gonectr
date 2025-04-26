package create

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func cloneOrUpdateGonerRepo() (string, error) {
	localDir, err := getRepoLocalDir(cacheDir, getGonerRepo())
	if err != nil {
		return "", err
	}
	err = cloneOrUpdateReop(localDir, getGonerRepo())
	return localDir, err
}

func listDirModule(dir string) ([]string, error) {
	var dirs []string
	return dirs, filepath.WalkDir(dir, func(p string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path: %v\n", err)
			return nil
		}
		if info.IsDir() {
			goModPath := filepath.Join(p, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				rel, _ := filepath.Rel(dir, p)
				dirs = append(dirs, rel)
				return filepath.SkipDir
			}
		}
		return nil
	})
}

func listTemplates() ([]string, error) {
	repo, err := cloneOrUpdateGonerRepo()
	if err != nil {
		return nil, err
	}
	return listDirModule(path.Join(repo, "examples"))
}

func listExamples() error {
	templates, err := listTemplates()
	if err != nil {
		return err
	}

	fmt.Println("Available templates:")
	for _, tpl := range templates {
		fmt.Printf("\t- %s\n", tpl)
	}
	return nil
}
