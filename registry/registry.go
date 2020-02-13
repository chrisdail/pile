package registry

// Registry operations that can be performed on a registry
type Registry interface {
	Prefix() string
	Contains(repository string, tag string) bool
}

// Config configuration for various registries
type Config struct {
	// Standard Docker Registry v2
	RegistryV2 DockerV2 `yaml:"registry_v2"`

	// Amazon ECR
	ECR AmazonECR
}

// ConfiguredRegistry determines which registry type is configured and returns Registry
func (config *Config) ConfiguredRegistry() Registry {
	if config.RegistryV2.URL != "" {
		return &config.RegistryV2
	} else if config.ECR.AccountID != "" {
		return &config.ECR
	}
	return nil
}

// RegistryPrefix return the prefix for this registry that is configured
func (config *Config) RegistryPrefix() string {
	registry := config.ConfiguredRegistry()
	if registry == nil {
		return ""
	}
	return registry.Prefix()
}
