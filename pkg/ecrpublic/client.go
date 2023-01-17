package ecrpublic

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/koki-develop/docker-tags/pkg/docker"
)

type Client struct {
	profile      string
	dockerClient *docker.Client
	httpClient   *http.Client
}

type Config struct {
	Profile string
}

func New(cfg *Config) *Client {
	return &Client{
		profile: cfg.Profile,
		dockerClient: docker.New(&docker.Config{
			APIURL:  "https://public.ecr.aws",
			AuthURL: "https://public.ecr.aws/v2/token",
		}),
		httpClient: new(http.Client),
	}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	tkn, err := cl.auth(name)
	if err != nil {
		return nil, err
	}

	tags, err := cl.listTags(name, tkn)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (cl *Client) auth(name string) (string, error) {
	cfg := &aws.Config{
		Region: aws.String("us-east-1"),
	}
	if cl.profile != "" {
		cfg.Credentials = credentials.NewSharedCredentials("", cl.profile)
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return "", err
	}

	svc := ecrpublic.New(sess)

	out, err := svc.GetAuthorizationToken(&ecrpublic.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}

	return *out.AuthorizationData.AuthorizationToken, nil
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (cl *Client) listTags(name, tkn string) ([]string, error) {
	req, err := cl.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	if tkn != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))
	}

	var tagsResp listTagsResponse
	if err := cl.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	return tagsResp.Tags, nil
}
