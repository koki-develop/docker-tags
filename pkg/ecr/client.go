package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Client struct {
	profile string
	region  string
}

type Config struct {
	Profile string
	Region  string
}

func New(cfg *Config) *Client {
	return &Client{
		profile: cfg.Profile,
		region:  cfg.Region,
	}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	cfg := &aws.Config{}
	if cl.region != "" {
		cfg.Region = &cl.region
	}
	if cl.profile != "" {
		cfg.Credentials = credentials.NewSharedCredentials("", cl.profile)
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
