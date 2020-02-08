package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chrisdail/pile/buildtools"
)

var tools = &buildtools.DockerBuildTools{}

var ErrorTestsFailed = errors.New("Tests Failed")

type Piler struct {
	Force     bool
	SkipPush  bool
	SkipTests bool
}

func (piler *Piler) Build(project *Project) (string, error) {
	if !project.CanBuild {
		return "", nil
	}

	if piler.Force {
		log.Println("Forcing rebuild due to --force flag")
	} else if project.GitVersion.Dirty {
		log.Println("Filesystem is dirty. Requires rebuilding")
	} // else if ExistsInRegistry() {}
	// log.Printf("Skipping build %s", image)
	// WriteManifest
	// return project.FullyQualifiedImage, nil
	//}

	if !piler.SkipTests && project.Config.Test.Target != "" {
		if err := piler.RunTests(project); err != nil {
			return "", err
		}
	}

	var buildImage = project.FullyQualifiedImage
	if piler.SkipPush {
		buildImage = project.Image
	}

	buildOptions := &buildtools.BuildOptions{
		Dir:       project.Dir,
		Image:     buildImage,
		Pull:      piler.Force,
		BuildArgs: project.Config.BuildArgs,
	}
	if err := tools.Build(buildOptions); err != nil {
		return buildImage, err
	}

	if piler.SkipPush {
		// Log
	} else {
		// Push image
	}

	// WriteManifest()
	return buildImage, nil
}

func (piler *Piler) RunTests(project *Project) error {
	testImage := fmt.Sprintf("%s-%s:%s", project.Config.Name, project.Config.Test.Target, project.Tag)

	buildOptions := &buildtools.BuildOptions{
		Dir:       project.Dir,
		Image:     testImage,
		Pull:      piler.Force,
		Target:    project.Config.Test.Target,
		BuildArgs: project.Config.BuildArgs,
	}
	if err := tools.Build(buildOptions); err != nil {
		return err
	}

	log.Printf("Running tests for %s using %s", project.Config.Name, testImage)
	rand.Seed(time.Now().UnixNano())
	containerName := fmt.Sprintf("pile-%s-%d", project.Config.Name, rand.Intn(100000))
	err := tools.Run(project.Dir, containerName, testImage)

	// TODO: Copy results

	// Remove the container. Intentionally ignore any errors
	tools.Rm(containerName)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("%s: %s", project.Config.Name, ErrorTestsFailed)
	}
	return nil
}
