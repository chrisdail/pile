package registry

import (
	"fmt"
)

const DefaultECRRegion = "us-east-1"

type RegistryConfig struct {
	// Standard Docker Registry v2
	RegistryV2 struct {
		URL string
	}
	// Amazon ECR
	ECR struct {
		AccountID string `yaml:"account_id"`
		Region    string
	}
}

func (config *RegistryConfig) RegistryPrefix() string {
	if config.RegistryV2.URL != "" {
		return fmt.Sprintf("%s/", config.RegistryV2.URL)
	} else if config.ECR.AccountID != "" {
		region := config.ECR.Region
		if region == "" {
			region = DefaultECRRegion
		}
		return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/", config.ECR.AccountID, region)
	}
	return ""
}
