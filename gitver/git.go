package gitver

import (
	"os/exec"
	"path/filepath"
	"strings"
)

var gitRoot = ""

// GitRootPath retrieves the absolute path of the root of this versioned git tree
func GitRootPath() (string, error) {
	var err error
	if gitRoot == "" {
		gitRoot, err = git("rev-parse", "--show-toplevel")
	}
	return gitRoot, err
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

func git(args ...string) (string, error) {
	output, err := exec.Command("git", args...).Output()
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
