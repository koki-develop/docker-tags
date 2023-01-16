package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Client struct {
	region string
}

type Config struct {
	Region string
}

func New(cfg *Config) *Client {
	return &Client{
		region: cfg.Region,
	}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	cfg := &aws.Config{}
	if cl.region != "" {
		cfg.Region = &cl.region
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	svc := ecr.New(sess)

	out, err := svc.DescribeImages(&ecr.DescribeImagesInput{
		RepositoryName: aws.String(name),
		Filter: &ecr.DescribeImagesFilter{
			TagStatus: aws.String(ecr.TagStatusTagged),
		},
	})
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for _, d := range out.ImageDetails {
		for _, t := range d.ImageTags {
			tags = append(tags, *t)
		}
	}

	return tags, nil
}
