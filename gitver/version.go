package gitver

import (
	"log"
	"os/user"
	"regexp"
	"strings"
	"text/template"
)

// DefaultTemplate Default template for formatting GitVersion using String()
const DefaultTemplate = "{{if .Dirty}}dirty-{{.User}}-{{end}}{{.Commits}}.{{.Hash}}"

// GitVersion version information about one or more git projects
type GitVersion struct {
	Commits string
	Hash    string
	Dirty   bool
	User    string
}

// FormatTemplate formats a GitVersion using a text/template string
func (ver *GitVersion) FormatTemplate(arg string) (string, error) {
	versionTemplate, err := template.New("version").Parse(arg)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if versionTemplate.Execute(&builder, ver) != nil {
		return "", err
	}
	return builder.String(), nil
}

func (ver *GitVersion) String() string {
	result, err := ver.FormatTemplate(DefaultTemplate)
	if err != nil {
		log.Println("Error formatting version template: %s", err)
		return ""
	}
	return result
}

// ForProjects computes the GitVersion for a set of projects relative to the git root
func ForProjects(projects []string) (*GitVersion, error) {
	paths, err := GitProjectPaths(projects)
	if err != nil {
		return nil, err
	}

	gitVersion := &GitVersion{}

	if commits, err := countCommits(paths); err == nil {
		gitVersion.Commits = commits
	}

	rev, err := headCommit(paths)
	if err != nil {
		return nil, err
	}
	if hash, err := revParseShort(rev); err == nil {
		gitVersion.Hash = hash
	}

	if dirty, err := checkIsDirty(paths); err == nil {
		gitVersion.Dirty = dirty
		gitVersion.User = currentUser()
	}

	return gitVersion, nil
}

func currentUser() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	alphaNumericPattern, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	return alphaNumericPattern.ReplaceAllString(user.Username, "")
}
