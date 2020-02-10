package registry

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

const DefaultECRRegion = "us-east-1"

type AmazonECR struct {
	AccountID string `yaml:"account_id"`
	Region    string
}

func (registry *AmazonECR) region() string {
	if registry.Region != "" {
		return registry.Region
	}
	return DefaultECRRegion
}

func (registry *AmazonECR) Prefix() string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/", registry.AccountID, registry.region())
}

func (registry *AmazonECR) Contains(repository string, tag string) bool {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(registry.region())))
	svc := ecr.New(sess)

	output, err := svc.DescribeImages(&ecr.DescribeImagesInput{
		RegistryId:     &registry.AccountID,
		RepositoryName: &repository,
		ImageIds:       []*ecr.ImageIdentifier{&ecr.ImageIdentifier{ImageTag: &tag}},
	})
	if err != nil {
		log.Println(fmt.Errorf("Error searching for image %s:%s in ECR: %v", repository, tag, err))
		return false
	}

	for _, detail := range output.ImageDetails {
		log.Printf("Image in registry %s:%s found, pushedAt %s, digest: %s",
			repository,
			tag,
			detail.ImagePushedAt,
			detail.ImageDigest,
		)
		return true
	}

	return false
}
