package gitver

import (
	"os/exec"
	"path/filepath"
	"strings"
)

var workingDir string
var gitRootCache = &cachedStringResponse{}
var gitBranchCache = &cachedStringResponse{}

// SetWorkingDir sets the working directory for running git commands
func SetWorkingDir(dir string) {
	workingDir = dir
}

// GitRootPath retrieves the absolute path of the root of this versioned git tree
func GitRootPath() (string, error) {
	return gitRootCache.cachedGit("rev-parse", "--show-toplevel")
}

// GitBranch gets the current git branch
func GitBranch() (string, error) {
	return gitBranchCache.cachedGit("branch", "--show-current")
}

// GitProjectPaths gives absolute paths given project paths relative to the git root
func GitProjectPaths(projects []string) ([]string, error) {
	rootPath, err := GitRootPath()
	if err != nil {
		return []string{}, err
	}

	paths := make([]string, len(projects))
	for i, project := range projects {
		paths[i] = filepath.Join(rootPath, project)
	}
	return paths, nil
}

func (cache *cachedStringResponse) cachedGit(args ...string) (string, error) {
	cache.Do(func() {
		cache.response, cache.err = git(args...)
	})
	return cache.response, cache.err
}

func git(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func countCommits(paths []string) (string, error) {
	commandArgs := []string{"rev-list", "HEAD", "--count", "--first-parent", "--"}
	commandArgs = append(commandArgs, paths...)
	return git(commandArgs...)
}

func headCommit(paths []string) (string, error) {
	commandArgs := []string{"rev-list", "-1", "HEAD", "--"}
	commandArgs = append(commandArgs, paths...)
	return git(commandArgs...)
}

func revParseShort(rev string) (string, error) {
	return git("rev-parse", "--short", rev)
}

func checkIsDirty(paths []string) (bool, error) {
	commandArgs := []string{"status", "--porcelain", "--"}
	commandArgs = append(commandArgs, paths...)
	status, err := git(commandArgs...)
	if err != nil {
		return false, err
	}

	return status != "", nil
}
