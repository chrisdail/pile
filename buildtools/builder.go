package buildtools

// BuildOptions options to a docker build command
type BuildOptions struct {
	Dir            string
	DockerfilePath string
	Image          string
	Target         string
	Pull           bool
	BuildArgs      map[string]string
}

// BuildTools interface for working with different types of builders
type BuildTools interface {
	Build(options *BuildOptions) error
	Push(image string) error

	Run(dir string, name string, image string) error
	Cp(src string, dst string) error
	Rm(name string) error
}
