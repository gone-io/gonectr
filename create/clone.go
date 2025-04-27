package create

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/gone-io/gonectr/utils"
	"net/url"
	"os"
	"path"
	"path/filepath"
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

func cloneOrUpdateReop(repoLocal string, repoURL string) (err error) {
	_, err = git.PlainClone(repoLocal, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			r, err := git.PlainOpen(repoLocal)
			if err != nil {
				return errors.Join(err, errors.New(fmt.Sprintf("open repo error, please delete `%s` and try again", repoLocal)))
			}
			w, err := r.Worktree()
			if err != nil {
				return errors.Join(err, errors.New(fmt.Sprintf("open repo error, please delete `%s` and try again", repoLocal)))
			}
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil {
				return errors.Join(err, errors.New("update repo error, please try again later"))
			}
		} else {
			return errors.Join(err, errors.New(fmt.Sprintf("git clone %s failed", repoURL)))
		}
	}
	return
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

		return "", false, fmt.Errorf("unsupported template name. Use gonectr create -ls to view supported template names;\n" +
			" alternatively, you can use the git repository address of a golang project.")
	}
	dir, _ = getRepoLocalDir(cacheDir, templateName)
	return dir, false, cloneOrUpdateReop(dir, templateName)
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
