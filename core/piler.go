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

// Piler options for performing a pile build operation
type Piler struct {
	Force     bool
	SkipPush  bool
	SkipTests bool
}

// BuildImage image data produced by a build
type BuildImage struct {
	Name                string `json:"name"`
	Repository          string `json:"repository"`
	Tag                 string `json:"tag"`
	FullyQualifiedImage string `json:"fully_qualified_image"`
}

// Build performs a build on a given project
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

	// Determine if we need to build or not
	if piler.Force {
		log.Println("Forcing rebuild due to --force flag")
	} else if project.GitVersion.Dirty {
		log.Println("Filesystem is dirty. Requires rebuilding")
	} else if project.Config.Registry.ConfiguredRegistry() != nil &&
		project.Config.Registry.ConfiguredRegistry().Contains(project.Repository, project.Tag) {

		log.Printf("Skipping build %s. Exists in registry", buildImage.FullyQualifiedImage)
		err := buildImage.WriteManifest(project.Dir)
		return buildImage, err
	}

	if !piler.SkipTests && project.Config.Test.Target != "" {
		if err := piler.RunTests(project); err != nil {
			return buildImage, err
		}
	}

	buildOptions := &buildtools.BuildOptions{
		Dir:            project.ContextDir(),
		DockerfilePath: filepath.Join(project.Dir, dockerfile),
		Image:          buildImage.FullyQualifiedImage,
		Pull:           piler.Force,
		BuildArgs:      project.Config.BuildArgs,
	}
	if err := tools.Build(buildOptions); err != nil {
		return buildImage, err
	}

	if piler.SkipPush {
		log.Printf("Skipping push %s\n", buildImage.FullyQualifiedImage)
	} else if project.Config.Registry.ConfiguredRegistry() != nil {
		log.Printf("Pushing image %s\n", buildImage.FullyQualifiedImage)
		if err := tools.Push(buildImage.FullyQualifiedImage); err != nil {
			return buildImage, err
		}
	}

	err := buildImage.WriteManifest(project.Dir)
	return buildImage, err
}

// WriteManifest writes out a descriptor about what was built
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

// RunTests runs tests for a project
func (piler *Piler) RunTests(project *Project) error {
	testImage := fmt.Sprintf("%s-%s:%s", project.Repository, project.Config.Test.Target, project.Tag)

	buildOptions := &buildtools.BuildOptions{
		Dir:            project.ContextDir(),
		DockerfilePath: filepath.Join(project.Dir, dockerfile),
		Image:          testImage,
		Pull:           piler.Force,
		Target:         project.Config.Test.Target,
		BuildArgs:      project.Config.BuildArgs,
	}
	if err := tools.Build(buildOptions); err != nil {
		return err
	}

	log.Printf("Running tests for %s using %s", project.Config.Name, testImage)
	rand.Seed(time.Now().UnixNano())
	containerName := fmt.Sprintf("pile-%s-%d", project.Config.Name, rand.Intn(100000))

	// Run test container in docker
	err := tools.Run(buildOptions.Dir, containerName, testImage)

	// Copy test results
	if project.Config.Test.CopyResults.SrcPath != "" && project.Config.Test.CopyResults.DstPath != "" {
		tools.Cp(
			fmt.Sprintf("%s:%s", containerName, project.Config.Test.CopyResults.SrcPath),
			filepath.Join(buildOptions.Dir, project.Config.Test.CopyResults.DstPath),
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
