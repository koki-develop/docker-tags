package ecr

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/koki-develop/docker-tags/pkg/util/awsutil"
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
	sess, err := awsutil.NewSession(&awsutil.SessionConfig{
		Region:  cl.extractRegionFromDomain(),
		Profile: cl.profile,
	})
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

	imgs := out.ImageDetails
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].ImagePushedAt.After(*imgs[j].ImagePushedAt)
	})

	tags := []string{}
	for _, img := range imgs {
		for _, t := range img.ImageTags {
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
