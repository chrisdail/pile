package core

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/chrisdail/pile/gitver"
)

type workspace struct {
	Dir string
}

// Workspace project workspace
var Workspace = &workspace{}

// SetWorkingDir sets the workspace working directory
func (ws *workspace) SetDir(dir string) error {
	if dir != "" {
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
