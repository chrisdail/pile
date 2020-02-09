package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

type BuildImage struct {
	Name                string `json:"name"`
	Repository          string `json:"repository"`
	Tag                 string `json:"tag"`
	FullyQualifiedImage string `json:"fully_qualified_image"`
}

func (piler *Piler) Build(project *Project) (*BuildImage, error) {
	if !project.CanBuild {
		return &BuildImage{}, nil
	}

	buildImage := &BuildImage{
		Name:                project.Config.Name,
		Repository:          project.Repository,
		Tag:                 project.Tag,
		FullyQualifiedImage: project.ImageWithRegistry,
	}

	if piler.SkipPush {
		buildImage.FullyQualifiedImage = project.Image
	}

	// Determine if we need to build or not
	if piler.Force {
		log.Println("Forcing rebuild due to --force flag")
	} else if project.GitVersion.Dirty {
		log.Println("Filesystem is dirty. Requires rebuilding")
	} // else if ExistsInRegistry() {}
	// log.Printf("Skipping build %s", image)
	// WriteManifest
	// buildImage.WriteManifest(project.Dir)
	// return buildImage.FullyQualifiedImage, nil
	//}

	if !piler.SkipTests && project.Config.Test.Target != "" {
		if err := piler.RunTests(project); err != nil {
			return buildImage, err
		}
	}

	buildOptions := &buildtools.BuildOptions{
		Dir:       project.Dir,
		Image:     buildImage.FullyQualifiedImage,
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

	err := buildImage.WriteManifest(project.Dir)
	return buildImage, err
}

func (image *BuildImage) WriteManifest(dir string) error {
	bytes, err := json.MarshalIndent(image, "", "    ")
	if err != nil {
		return err
	}

	buildDir := filepath.Join(dir, "build")
	os.MkdirAll(buildDir, os.ModePerm)
	descriptorPath := filepath.Join(buildDir, "pile-image.json")
	err = ioutil.WriteFile(descriptorPath, bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (piler *Piler) RunTests(project *Project) error {
	testImage := fmt.Sprintf("%s-%s:%s", project.Repository, project.Config.Test.Target, project.Tag)

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

	// Run test container in docker
	err := tools.Run(project.Dir, containerName, testImage)

	// Copy test results
	if project.Config.Test.CopyResults.SrcPath != "" && project.Config.Test.CopyResults.DstPath != "" {
		tools.Cp(
			fmt.Sprintf("%s:%s", containerName, project.Config.Test.CopyResults.SrcPath),
			filepath.Join(project.Dir, project.Config.Test.CopyResults.DstPath),
		)
	}

	// Remove the container. Intentionally ignore any errors
	tools.Rm(containerName)

	if err != nil {
		log.Println(err)
		return fmt.Errorf("%s: %s", project.Config.Name, ErrorTestsFailed)
	}
	return nil
}
