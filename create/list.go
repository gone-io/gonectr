package create

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func cloneOrUpdateGonerRepo() (string, error) {
	localDir, err := getRepoLocalDir(cacheDir, getGonerRepo())
	if err != nil {
		return "", err
	}
	err = cloneOrUpdateRepo(localDir, getGonerRepo())
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
	localDir, err := getRepoLocalDir(cacheDir, getGonerRepo())
	if err != nil {
		return err
	}

	templates, err := listTemplates()
	if err != nil {
		return err
	}

	fmt.Println("Available templates:")

	for _, tpl := range templates {
		readMePath := path.Join(localDir, "examples", tpl, "README.md")
		var desc string
		if stat, err := os.Stat(readMePath); err == nil && !stat.IsDir() {
			desc = getReadmeDesc(readMePath)
		}

		fmt.Printf("  - %-20s\t%s\n", tpl, desc)
	}
	return nil
}

func getReadmeDesc(filename string) string {
	content, err := os.ReadFile(filename)
	if err == nil {
		re := regexp.MustCompile(`(?s)\[//\]: # \(desc[:：]\s*(.*?)\s*\)`)
		matches := re.FindStringSubmatch(string(content))
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
