package ecrpublic

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/aws/aws-sdk-go/service/ecrpublic/ecrpubliciface"
	"github.com/koki-develop/docker-tags/pkg/util/awsutil"
	"github.com/koki-develop/docker-tags/pkg/util/dockerutil"
)

type Client struct {
	ecrpublicAPI ecrpubliciface.ECRPublicAPI

	dockerClient *dockerutil.Client
	httpClient   *http.Client
}

type Config struct {
	Profile string
}

func New(cfg *Config) (*Client, error) {
	sess, err := awsutil.NewSession(&awsutil.SessionConfig{
		Profile: cfg.Profile,
		Region:  "us-east-1",
	})
	if err != nil {
		return nil, err
	}

	svc := ecrpublic.New(sess)

	return &Client{
		ecrpublicAPI: svc,
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  "https://public.ecr.aws",
			AuthURL: "https://public.ecr.aws/v2/token",
		}),
		httpClient: new(http.Client),
	}, nil
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
	out, err := cl.ecrpublicAPI.GetAuthorizationToken(&ecrpublic.GetAuthorizationTokenInput{})
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
