package registry

import "fmt"

type RegistryDockerV2 struct {
	URL string
}

func (registry *RegistryDockerV2) Prefix() string {
	return fmt.Sprintf("%s/", registry.URL)
}

func (registry *RegistryDockerV2) Contains(repository string, tag string) (bool, error) {
	return false, nil
}
