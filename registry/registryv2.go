package registry

import (
	"fmt"
	"log"
	"net/url"

	"github.com/heroku/docker-registry-client/registry"
)

type RegistryDockerV2 struct {
	URL      string
	Insecure bool
}

func (docker *RegistryDockerV2) Prefix() string {
	parsed, err := url.Parse(docker.URL)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("%s/", parsed.Host)
}

func (docker *RegistryDockerV2) Contains(repository string, tag string) bool {
	var reg *registry.Registry
	var err error
	if docker.Insecure {
		reg, err = registry.NewInsecure(docker.URL, "", "")
	} else {
		reg, err = registry.New(docker.URL, "", "")
	}
	if err != nil {
		log.Println(err)
		return false
	}

	tags, err := reg.Tags(repository)
	if err != nil {
		log.Println(err)
		return false
	}

	for _, found := range tags {
		log.Println(found)
		if tag == found {
			return true
		}
	}

	return false
}
