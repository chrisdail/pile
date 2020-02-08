package core

import (
	"log"

	"github.com/chrisdail/pile/buildtools"
)

var tools = &buildtools.DockerBuildTools{}

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

	if !piler.SkipTests {
		// Run Tests
	}

	// TODO: Do build

	var buildImage = project.FullyQualifiedImage
	if piler.SkipPush {
		buildImage = project.Image
	}

	log.Printf("Building %s\n", buildImage)

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
