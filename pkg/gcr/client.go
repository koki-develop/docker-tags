package gcr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/koki-develop/docker-tags/pkg/docker"
	"golang.org/x/oauth2/google"
)

type Client struct {
	dockerClient *docker.Client
	httpClient   *http.Client
}

func New() *Client {
	return &Client{
		dockerClient: docker.New(&docker.Config{
			APIURL:  "https://gcr.io",
			AuthURL: "https://gcr.io/v2/token",
		}),
		httpClient: new(http.Client),
	}
}

type listTagsResponse struct {
	Manifest map[string]manifest `json:"manifest"`
}

type manifest struct {
	Tag []string `json:"tag"`
}

func (cl *Client) ListTags(name string) ([]string, error) {
	tkn, _ := cl.auth(name)

	tags, err := cl.listTags(name, tkn)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (cl *Client) auth(name string) (string, error) {
	ctx := context.Background()
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", err
	}

	tkn, err := cred.TokenSource.Token()
	if err != nil {
		return "", err
	}

	req, err := cl.dockerClient.NewAuthRequest(name)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("_token", tkn.AccessToken)

	resp, err := cl.dockerClient.DoAuthRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
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

	tags := []string{}
	for _, m := range tagsResp.Manifest {
		tags = append(tags, m.Tag...)
	}

	return tags, nil
}
