package buildtools

type BuildOptions struct {
	Dir       string
	File      string
	Image     string
	Target    string
	Pull      bool
	BuildArgs map[string]string
}

type BuildTools interface {
	Build(options *BuildOptions) error
	Push(image string) error

	Run(dir string, name string, image string) error
	Cp(src string, dst string) error
	Rm(name string) error
}
