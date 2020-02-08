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
}
