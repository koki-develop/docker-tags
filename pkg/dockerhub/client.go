package dockerhub

import (
	"fmt"
	"net/http"

	"github.com/koki-develop/docker-tags/pkg/docker"
)

type Client struct {
	dockerClient *docker.Client
	httpClient   *http.Client
}

func New() *Client {
	return &Client{
		dockerClient: docker.New(&docker.Config{
			APIURL:  "https://registry.hub.docker.com",
			AuthURL: "https://auth.docker.io/token",
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

	// reverse
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}

	return tags, nil
}

func (cl *Client) auth(name string) (string, error) {
	req, err := cl.dockerClient.NewAuthRequest(name)
	if err != nil {
		return "", err
	}

	resp, err := cl.dockerClient.DoAuthRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (cl *Client) listTags(name, tkn string) ([]string, error) {
	req, err := cl.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))

	var tagsResp listTagsResponse
	if err := cl.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	return tagsResp.Tags, nil
}
