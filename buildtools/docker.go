package buildtools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type DockerBuildTools struct{}

func (*DockerBuildTools) Build(options *BuildOptions) error {
	log.Printf("Building %s\n", options.Image)

	args := []string{"build", "."}
	if options.File != "" {
		args = append(args,
			"-f",
			filepath.Join(options.Dir, options.File),
		)
	}
	if options.Target != "" {
		args = append(args, "--target", options.Target)
	}
	if options.Pull {
		args = append(args, "--pull")
	}
	for key, value := range options.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	args = append(args, "-t", options.Image)
	return docker(options.Dir, args...)
}

func (*DockerBuildTools) Push(image string) error {
	return docker("", "push", image)
}

func (*DockerBuildTools) Run(dir string, name string, image string) error {
	return docker(dir, "run", "-t", "--name", name, image)
}

func (*DockerBuildTools) Cp(src string, dst string) error {
	return docker("", "cp", src, dst)
}

func (*DockerBuildTools) Rm(name string) error {
	return docker("", "rm", "-f", name)
}

func docker(dir string, args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing docker %s: %w", args, err)
	}
	return nil
}
