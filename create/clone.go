package create

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gone-io/gonectl/utils"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func getRepoLocalDir(cacheDir string, repoURL string) (string, error) {
	Url, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, Url.Host, Url.Path)
	if filepath.Ext(dir) == ".git" {
		dir = strings.TrimSuffix(dir, ".git")
	}
	return dir, nil
}

func cloneOrUpdateRepo(repoLocal string, repoURL string) (err error) {
	// 检查本地仓库是否存在
	_, err = os.Stat(repoLocal)
	var repo *git.Repository

	if os.IsNotExist(err) {
		// 如果本地仓库不存在，克隆仓库
		fmt.Printf("Cloning %s to %s\n", repoURL, repoLocal)
		// 确保父目录存在
		if err = os.MkdirAll(filepath.Dir(repoLocal), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}

		// 克隆仓库
		repo, err = git.PlainClone(repoLocal, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %v", err)
		}
	} else {
		// 如果本地仓库已存在，打开仓库
		repo, err = git.PlainOpen(repoLocal)
		if err != nil {
			return fmt.Errorf("failed to open repository: %v", err)
		}

		// 获取远程仓库
		remote, err := repo.Remote("origin")
		if err != nil {
			return fmt.Errorf("failed to get remote: %v", err)
		}

		// 执行fetch更新代码
		fmt.Printf("Fetching updates for %s\n", repoURL)
		err = remote.Fetch(&git.FetchOptions{
			Progress: os.Stdout,
			Tags:     git.AllTags,
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			_ = fmt.Errorf("failed to fetch updates: %v", err)
		}
	}

	// 获取所有标签
	tagRefs, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("failed to get tags: %v", err)
	}

	// 查找最新的语义化版本标签
	var latestTag string
	var latestTagCommit *object.Commit

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		tagName := tagRef.Name().Short()
		// 简单的语义化版本检查，可以根据需要使用更复杂的规则
		if suc, err := regexp.MatchString(`^v[0-9]+\.[0-9]+\.[0-9]+$`, tagName); err == nil && suc {
			tagObj, err := repo.TagObject(tagRef.Hash())
			if err == nil {
				// 这是一个标注标签
				commit, err := tagObj.Commit()
				if err != nil {
					return nil // 跳过无法获取提交的标签
				}
				if latestTag == "" || strings.Compare(tagName, latestTag) > 0 {
					latestTag = tagName
					latestTagCommit = commit
				}
			} else if errors.Is(err, plumbing.ErrObjectNotFound) {
				// 这是一个轻量级标签
				commit, err := repo.CommitObject(tagRef.Hash())
				if err != nil {
					return nil // 跳过无法获取提交的标签
				}
				if latestTag == "" || strings.Compare(tagName, latestTag) > 0 {
					latestTag = tagName
					latestTagCommit = commit
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to iterate tags: %v", err)
	}

	// 如果找到了语义化版本标签，切换到该版本
	if latestTag != "" && latestTagCommit != nil {
		fmt.Printf("Checking out latest version: %s\n", latestTag)

		// 获取工作树
		worktree, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %v", err)
		}

		// 切换到最新版本
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: latestTagCommit.Hash,
		})
		if err != nil {
			return fmt.Errorf("failed to checkout tag %s: %v", latestTag, err)
		}
		fmt.Printf("Successfully checked out version %s\n", latestTag)
	} else {
		fmt.Println("No semantic version tags found, using latest commit")
	}
	return nil
}

func getGonerRepo() string {
	repUrl := "https://github.com/gone-io/goner.git"
	if utils.IsInChina() {
		repUrl = "https://gitee.com/gone-io/goner.git"
	}
	return repUrl
}

func getTemplateCodeDir(templateName string, cacheDir string) (dir string, isExample bool, err error) {
	parse, err := url.Parse(templateName)
	if err != nil {
		return "", false, err
	}
	if parse.Scheme == "" && parse.Host == "" {
		templates, err := listTemplates()
		if err != nil {
			return "", false, err
		}

		for _, tpl := range templates {
			if tpl == templateName {
				dir, _ = getRepoLocalDir(cacheDir, getGonerRepo())
				return path.Join(dir, "examples", templateName), true, nil
			}
		}

		return "", false, fmt.Errorf("unsupported template name. Use gonectl create -ls to view supported template names;\n" +
			" alternatively, you can use the git repository address of a golang project.")
	}
	dir, _ = getRepoLocalDir(cacheDir, templateName)
	return dir, false, cloneOrUpdateRepo(dir, templateName)
}

func checkDirIsGoModuleAndGetModuleName(dir string) (string, error) {
	modFile := path.Join(dir, "go.mod")
	stat, err := os.Stat(modFile)
	if err != nil {
		return "", err
	}
	if stat.IsDir() {
		return "", errors.New("go.mod is a dir")
	}
	return utils.ParseModuleName(modFile)
}
