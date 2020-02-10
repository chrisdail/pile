package core

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/chrisdail/pile/gitver"
)

type workspace struct {
	Project
}

// Workspace project workspace
var Workspace = &workspace{}

// SetWorkingDir sets the workspace working directory
func (ws *workspace) SetDir(dir string) error {
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return err
		}

		ws.Dir = dir
	} else {
		var err error
		ws.Dir, err = gitver.GitRootPath()
		if err != nil {
			return err
		}
	}
	return nil
}

// ProjectPaths get paths for corresponding project names
func (ws *workspace) ProjectPaths(projects []string) []string {
	if ws.Dir == "" {
		log.Fatalln(errors.New("No workspace found"))
	}

	paths := make([]string, len(projects))
	for i, project := range projects {
		paths[i] = filepath.Join(ws.Dir, project)
	}
	return paths
}

// DiscoverProjectPaths walks the workspace tree, searching for project paths
func (ws *workspace) DiscoverProjectPaths() ([]string, error) {
	var paths []string
	err := filepath.Walk(ws.Dir, func(path string, info os.FileInfo, err error) error {
		if info.Name() == pileConfigName {
			paths = append(paths, filepath.Dir(path))
		}
		return nil
	})
	return paths, err
}

func (ws *workspace) ProjectsFromArgs(args []string) ([]Project, error) {
	var (
		paths []string
		err   error
	)
	if len(args) == 0 {
		paths, err = ws.DiscoverProjectPaths()
		if err != nil {
			return nil, err
		}
	} else {
		paths = ws.ProjectPaths(args)
	}

	return ws.loadProjects(paths)
}

func (ws *workspace) loadProjects(paths []string) ([]Project, error) {
	err := ws.Load(&ProjectConfig{})
	if err != nil {
		return nil, err
	}

	// Remove the name from the workspace project as it is not a default
	ws.Config.Name = ""

	projects := make([]Project, len(paths))
	for i, path := range paths {
		if ws.Dir == path {
			projects[i] = ws.Project
		} else {
			projects[i] = Project{Dir: path}
			err := projects[i].Load(&ws.Config)
			if err != nil {
				return projects, err
			}
		}
	}
	return projects, nil
}
