// Package git provides Git repository integration for file status tracking.
package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Repository struct {
	repoRoot     string
	repo         *git.Repository
	cachedStatus git.Status
	statusLoaded bool
}

func NewRepository(path string) (*Repository, error) {
	root, err := findGitRoot(path)
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(root)
	if err != nil {
		return nil, err
	}
	return &Repository{repoRoot: root, repo: repo}, nil
}

func (g *Repository) loadStatus() error {
	if g.statusLoaded {
		return nil
	}
	worktree, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("getting worktree status: %w", err)
	}
	g.cachedStatus = status
	g.statusLoaded = true
	return nil
}

func (g *Repository) GetStatus(filePath string) string {
	if err := g.loadStatus(); err != nil {
		return ""
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	}

	relPath, err := filepath.Rel(g.repoRoot, absPath)
	if err != nil {
		return ""
	}

	fileStatus := g.cachedStatus.File(relPath)

	if fileStatus.Worktree == git.Untracked {
		return "?"
	}

	if fileStatus.Worktree == git.Unmodified && fileStatus.Staging == git.Unmodified {
		return ""
	}

	var statusStr string
	if fileStatus.Staging != git.Unmodified {
		switch fileStatus.Staging {
		case git.Added:
			statusStr += "A"
		case git.Modified:
			statusStr += "M"
		case git.Deleted:
			statusStr += "D"
		case git.Renamed:
			statusStr += "R"
		case git.Copied:
			statusStr += "C"
		}
	} else {
		statusStr += " "
	}

	if fileStatus.Worktree != git.Unmodified {
		switch fileStatus.Worktree {
		case git.Modified:
			statusStr += "M"
		case git.Deleted:
			statusStr += "D"
		case git.Added:
			statusStr += "A"
		}
	}

	if statusStr == " " || statusStr == "" {
		return ""
	}

	return strings.TrimSpace(statusStr)
}

func findGitRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not a git repository")
		}
		dir = parent
	}
}
