package create

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Test_cloneOrUpdateRepo(t *testing.T) {
	// 创建临时目录作为测试仓库
	testRepoDir, err := os.MkdirTemp("", "test-repo-*")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(testRepoDir)

	// 创建临时目录作为本地仓库目标
	testLocalDir, err := os.MkdirTemp("", "test-local-*")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(testLocalDir)

	// 初始化测试仓库
	repo, err := git.PlainInit(testRepoDir, false)
	if err != nil {
		t.Fatalf("无法初始化Git仓库: %v", err)
	}

	// 创建工作树
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("无法获取工作树: %v", err)
	}

	// 创建一个测试文件
	testFilePath := filepath.Join(testRepoDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("initial content"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	// 添加文件到暂存区
	_, err = worktree.Add("test.txt")
	if err != nil {
		t.Fatalf("无法添加文件到暂存区: %v", err)
	}

	// 提交更改
	commit, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("无法提交更改: %v", err)
	}

	// 创建一个语义化版本标签
	_, err = repo.CreateTag("v1.0.0", commit, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Message: "Version 1.0.0",
	})
	if err != nil {
		t.Fatalf("无法创建标签: %v", err)
	}

	// 测试场景1: 克隆新仓库
	t.Run("克隆新仓库", func(t *testing.T) {
		localRepoPath := filepath.Join(testLocalDir, "repo1")
		err := cloneOrUpdateRepo(localRepoPath, "file://"+testRepoDir)
		if err != nil {
			t.Fatalf("克隆仓库失败: %v", err)
		}

		// 验证仓库是否被正确克隆
		_, err = git.PlainOpen(localRepoPath)
		if err != nil {
			t.Fatalf("无法打开克隆的仓库: %v", err)
		}

		// 验证文件是否存在
		clonedFilePath := filepath.Join(localRepoPath, "test.txt")
		if _, err := os.Stat(clonedFilePath); os.IsNotExist(err) {
			t.Fatalf("克隆的仓库中缺少测试文件")
		}
	})

	// 为更新测试准备: 修改原始仓库并创建新标签
	if err := os.WriteFile(testFilePath, []byte("updated content"), 0644); err != nil {
		t.Fatalf("无法更新测试文件: %v", err)
	}

	_, err = worktree.Add("test.txt")
	if err != nil {
		t.Fatalf("无法添加更新的文件到暂存区: %v", err)
	}

	newCommit, err := worktree.Commit("Update test file", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("无法提交更新: %v", err)
	}

	// 创建新的语义化版本标签
	_, err = repo.CreateTag("v1.1.0", newCommit, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
		Message: "Version 1.1.0",
	})
	if err != nil {
		t.Fatalf("无法创建新标签: %v", err)
	}

	// 测试场景2: 更新已存在的仓库
	t.Run("更新已存在的仓库", func(t *testing.T) {
		localRepoPath := filepath.Join(testLocalDir, "repo2")

		// 首先克隆仓库
		err := cloneOrUpdateRepo(localRepoPath, "file://"+testRepoDir)
		if err != nil {
			t.Fatalf("初始克隆仓库失败: %v", err)
		}

		// 再次调用函数更新仓库
		err = cloneOrUpdateRepo(localRepoPath, "file://"+testRepoDir)
		if err != nil {
			t.Fatalf("更新仓库失败: %v", err)
		}

		// 验证仓库是否被更新到最新版本
		repo, err := git.PlainOpen(localRepoPath)
		if err != nil {
			t.Fatalf("无法打开更新的仓库: %v", err)
		}

		// 获取HEAD引用
		head, err := repo.Head()
		if err != nil {
			t.Fatalf("无法获取HEAD引用: %v", err)
		}

		// 获取最新标签的引用
		tagRef, err := repo.Tag("v1.1.0")
		if err != nil {
			t.Fatalf("无法获取标签引用: %v", err)
		}

		// 获取标签对象
		tagObj, err := repo.TagObject(tagRef.Hash())
		if err != nil {
			// 如果是轻量级标签，直接比较引用哈希
			if head.Hash().String() != tagRef.Hash().String() {
				t.Fatalf("仓库未更新到最新标签: 期望 %s, 实际 %s", tagRef.Hash().String(), head.Hash().String())
			}
		} else {
			// 如果是标注标签，获取提交对象并比较
			commit, err := tagObj.Commit()
			if err != nil {
				t.Fatalf("无法获取标签提交: %v", err)
			}
			if head.Hash().String() != commit.Hash.String() {
				t.Fatalf("仓库未更新到最新标签: 期望 %s, 实际 %s", commit.Hash.String(), head.Hash().String())
			}
		}
	})

	// 测试场景3: 处理没有语义化版本标签的仓库
	t.Run("处理没有语义化版本标签的仓库", func(t *testing.T) {
		// 创建一个新的测试仓库，不添加语义化版本标签
		noTagRepoDir, err := os.MkdirTemp("", "test-no-tag-repo-*")
		if err != nil {
			t.Fatalf("无法创建临时目录: %v", err)
		}
		defer os.RemoveAll(noTagRepoDir)

		// 初始化仓库
		noTagRepo, err := git.PlainInit(noTagRepoDir, false)
		if err != nil {
			t.Fatalf("无法初始化Git仓库: %v", err)
		}

		// 创建工作树
		noTagWorktree, err := noTagRepo.Worktree()
		if err != nil {
			t.Fatalf("无法获取工作树: %v", err)
		}

		// 创建测试文件
		noTagFilePath := filepath.Join(noTagRepoDir, "test.txt")
		if err := os.WriteFile(noTagFilePath, []byte("no tag content"), 0644); err != nil {
			t.Fatalf("无法创建测试文件: %v", err)
		}

		// 添加并提交
		_, err = noTagWorktree.Add("test.txt")
		if err != nil {
			t.Fatalf("无法添加文件到暂存区: %v", err)
		}

		noTagCommit, err := noTagWorktree.Commit("Initial commit without tag", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Test User",
				Email: "test@example.com",
				When:  time.Now(),
			},
		})
		if err != nil {
			t.Fatalf("无法提交更改: %v", err)
		}

		// 创建一个非语义化版本标签
		_, err = noTagRepo.CreateTag("release-2023", noTagCommit, &git.CreateTagOptions{
			Tagger: &object.Signature{
				Name:  "Test User",
				Email: "test@example.com",
				When:  time.Now(),
			},
			Message: "Release 2023",
		})
		if err != nil {
			t.Fatalf("无法创建非语义化标签: %v", err)
		}

		// 克隆仓库
		localNoTagPath := filepath.Join(testLocalDir, "repo-no-tag")
		err = cloneOrUpdateRepo(localNoTagPath, "file://"+noTagRepoDir)
		if err != nil {
			t.Fatalf("克隆无标签仓库失败: %v", err)
		}

		// 验证仓库是否被正确克隆
		clonedRepo, err := git.PlainOpen(localNoTagPath)
		if err != nil {
			t.Fatalf("无法打开克隆的仓库: %v", err)
		}

		// 获取HEAD引用
		head, err := clonedRepo.Head()
		if err != nil {
			t.Fatalf("无法获取HEAD引用: %v", err)
		}

		// 验证HEAD是否指向最新提交
		if head.Hash().String() != noTagCommit.String() {
			t.Fatalf("仓库未克隆到最新提交: 期望 %s, 实际 %s", noTagCommit.String(), head.Hash().String())
		}
	})
}
