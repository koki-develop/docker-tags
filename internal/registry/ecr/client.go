package ecr

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/koki-develop/docker-tags/internal/util/awsutil"
)

type Client struct {
	ecrAPI ecriface.ECRAPI
}

type Config struct {
	Profile string
	Domain  string
}

func New(cfg *Config) (*Client, error) {
	sess, err := awsutil.NewSession(&awsutil.SessionConfig{
		Profile: cfg.Profile,
		Region:  extractRegionFromDomain(cfg.Domain),
	})
	if err != nil {
		return nil, err
	}

	svc := ecr.New(sess)

	return &Client{ecrAPI: svc}, nil
}

func extractRegionFromDomain(domain string) string {
	// <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com
	ds := strings.Split(domain, ".")
	return ds[len(ds)-3]
}

func (cl *Client) ListTags(name string) ([]string, error) {
	out, err := cl.ecrAPI.DescribeImages(&ecr.DescribeImagesInput{
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
