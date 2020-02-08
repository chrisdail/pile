package buildtools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type DockerBuildTools struct{}

func (*DockerBuildTools) Build(options *BuildOptions) error {
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

func docker(dir string, args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error executing docker %s: %w", args, err)
	}
	return nil
}
