package ecr

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Client struct {
	profile string
	domain  string
}

type Config struct {
	Profile string
	Domain  string
}

func New(cfg *Config) *Client {
	return &Client{
		profile: cfg.Profile,
		domain:  cfg.Domain,
	}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	cfg := &aws.Config{
		Region: aws.String(cl.extractRegionFromDomain()),
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

func (cl *Client) extractRegionFromDomain() string {
	// <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPOSITORY_NAME>
	ds := strings.Split(cl.domain, ".")
	return ds[len(ds)-3]
}
