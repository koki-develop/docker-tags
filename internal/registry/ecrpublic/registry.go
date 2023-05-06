package ecrpublic

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/aws/aws-sdk-go/service/ecrpublic/ecrpubliciface"
	"github.com/koki-develop/docker-tags/internal/util/awsutil"
	"github.com/koki-develop/docker-tags/internal/util/dockerutil"
)

type Registry struct {
	ecrpublicAPI ecrpubliciface.ECRPublicAPI

	dockerClient *dockerutil.Client
	httpClient   *http.Client
}

type Config struct {
	Profile string
}

func New(cfg *Config) (*Registry, error) {
	sess, err := awsutil.NewSession(&awsutil.SessionConfig{
		Profile: cfg.Profile,
		Region:  "us-east-1",
	})
	if err != nil {
		return nil, err
	}

	svc := ecrpublic.New(sess)

	return &Registry{
		ecrpublicAPI: svc,
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  "https://public.ecr.aws",
			AuthURL: "https://public.ecr.aws/v2/token",
		}),
		httpClient: new(http.Client),
	}, nil
}

func (r *Registry) ListTags(name string) ([]string, error) {
	tkn, err := r.auth(name)
	if err != nil {
		return nil, err
	}

	tags, err := r.listTags(name, tkn)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Registry) auth(name string) (string, error) {
	out, err := r.ecrpublicAPI.GetAuthorizationToken(&ecrpublic.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}

	return *out.AuthorizationData.AuthorizationToken, nil
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (r *Registry) listTags(name, tkn string) ([]string, error) {
	req, err := r.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	if tkn != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))
	}

	var tagsResp listTagsResponse
	if err := r.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	return tagsResp.Tags, nil
}
