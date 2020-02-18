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

var sanitizedUserCache = &cachedStringResponse{}

// GitVersion version information about one or more git projects
type GitVersion struct {
	Branch  string
	Commits string
	Hash    string
	Dirty   bool
	User    string
}

// FormatTemplate formats a GitVersion using a text/template string
func (ver *GitVersion) FormatTemplate(templateString string) (string, error) {
	if templateString == "" {
		templateString = DefaultTemplate
	}
	versionTemplate, err := template.New("Version Template").Parse(templateString)
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
		log.Fatalf("Error formatting version template: %s", err)
		return ""
	}
	return result
}

// New creates a new GitVersion for the specified paths
func New(paths []string) (*GitVersion, error) {
	var version = &GitVersion{}
	if err := version.forPaths(paths); err != nil {
		return nil, err
	}
	return version, nil
}

// ForProjects computes the GitVersion for a set of projects relative to the git root
func (ver *GitVersion) forPaths(paths []string) error {
	var err error

	// Ignore errors on git branch. Could be a detached head
	ver.Branch, _ = GitBranch()

	if ver.Commits, err = countCommits(paths); err != nil {
		return err
	}

	rev, err := headCommit(paths)
	if err != nil {
		return err
	}
	if rev == "" {
		ver.Hash = "untracked"
	} else if ver.Hash, err = revParseShort(rev); err != nil {
		return err
	}

	if ver.Dirty, err = checkIsDirty(paths); err != nil {
		return err
	}

	if ver.User, err = currentUser(); err != nil {
		return err
	}
	return nil
}

func currentUser() (string, error) {
	sanitizedUserCache.Do(func() {
		var currentUser *user.User
		currentUser, sanitizedUserCache.err = user.Current()
		if sanitizedUserCache.err != nil {
			return
		}

		var alphaNumericPattern *regexp.Regexp
		alphaNumericPattern, sanitizedUserCache.err = regexp.Compile("[^a-zA-Z0-9]+")
		if sanitizedUserCache.err != nil {
			return
		}

		sanitizedUserCache.response = alphaNumericPattern.ReplaceAllString(currentUser.Username, "")
		sanitizedUserCache.err = nil
	})
	return sanitizedUserCache.response, sanitizedUserCache.err
}
