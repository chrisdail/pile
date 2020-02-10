package registry

import "fmt"

const DefaultECRRegion = "us-east-1"

type AmazonECR struct {
	AccountID string `yaml:"account_id"`
	Region    string
}

func (registry *AmazonECR) Prefix() string {
	region := registry.Region
	if region == "" {
		region = DefaultECRRegion
	}
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/", registry.AccountID, region)
}

func (registry *AmazonECR) Contains(repository string, tag string) (bool, error) {
	return false, nil
}
