package registry

type Registry interface {
	Prefix() string
	Contains(repository string, tag string) bool
}

type RegistryConfig struct {
	// Standard Docker Registry v2
	RegistryV2 RegistryDockerV2 `yaml:"registry_v2"`

	// Amazon ECR
	ECR AmazonECR
}

func (config *RegistryConfig) ConfiguredRegistry() Registry {
	if config.RegistryV2.URL != "" {
		return &config.RegistryV2
	} else if config.ECR.AccountID != "" {
		return &config.ECR
	}
	return nil
}

func (config *RegistryConfig) RegistryPrefix() string {
	registry := config.ConfiguredRegistry()
	if registry == nil {
		return ""
	}
	return registry.Prefix()
}
