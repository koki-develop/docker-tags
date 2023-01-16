package ecr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	sess, err := session.NewSession(&aws.Config{})
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
