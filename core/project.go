package core

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chrisdail/pile/gitver"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

const PileConfigName = "pile.yml"

// #
// # Containers can have a metadata config file called 'buildc.yml' in the root of their tree containing:
// #
// # name: alt-name                   // Alternative name for this image. Defaults to the directory name
// # group: apps                      // Alternative sub-group for this image. Defaults to 'development'
// # depends_on:
// #   - other/project                // Optionally include other projects in the version.sh hash
// # args:
// #   KEY: VALUE                     // A map of arguments for the docker build: `--build-arg`
// # test:                            // Optional testing section
// #  target: test                    // Target in the Dockerfile for the test runner
// #  copy_results:                   // Copy test results from the container to the local filesystem (docker cp)
// #    src_path: /app/build/.        // Location to copy files from in the container
// #    dest_path: build              // Location to copy files to relative to the project directory
// # alt_container:                   // Root level options for an alternate container (specified via container parameter init)

type ProjectConfig struct {
	Name            string
	ImagePrefix     string            `yaml:"image_prefix"`
	VersionPrefix   string            `yaml:"version_prefix"`
	VersionTemplate string            `yaml:"version_template"`
	DependsOn       []string          `yaml:"depends_on"`
	BuildArgs       map[string]string `yaml:"build_args"`
	Test            struct {
		Target      string
		CopyResults struct {
			SrcPath string `yaml:"src_path"`
			DstPath string `yaml:"dst_path"`
		} `yaml:"copy_results"`
	}
	Registry struct {
		ECR struct {
			AccountID string `yaml:"account_id"`
			Region    string
		}
	}
}

type Project struct {
	Dir    string
	Config ProjectConfig
}

func (project *Project) Name() string {
	if project.Config.Name != "" {
		return project.Config.Name
	}
	return filepath.Dir(project.Dir)
}

func (project *Project) Load() {
	configPath := filepath.Join(project.Dir, PileConfigName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file does not exist: %s", configPath)
		return
	}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("Error reading config file %s: %s\n", configPath, err)
		return
	}

	err = yaml.Unmarshal(configFile, &project.Config)
	if err != nil {
		log.Printf("Error parsing YAML: %s\n", err)
		return
	}
}

func (project *Project) LoadWithDefaults(defaults *ProjectConfig) {
	project.Load()
	mergo.Merge(&project.Config, defaults)
}

func (project *Project) Version() (string, error) {
	// TODO: Handle dependencies
	version, err := gitver.New([]string{project.Dir})
	if err != nil {
		return "", err
	}

	renderedVersion, err := version.FormatTemplate(project.Config.VersionTemplate)
	if err != nil {
		return "", err
	}
	if project.Config.VersionPrefix != "" {
		return project.Config.VersionPrefix + renderedVersion, nil
	}
	return renderedVersion, nil
}
