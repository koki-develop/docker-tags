package dockerhub

import (
	"fmt"
	"net/http"

	"github.com/koki-develop/docker-tags/pkg/docker"
)

type Client struct {
	token        string
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
	if err := cl.auth(name); err != nil {
		return nil, err
	}

	tags, err := cl.listTags(name)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (cl *Client) auth(name string) error {
	req, err := cl.dockerClient.NewAuthRequest(name)
	if err != nil {
		return err
	}

	resp, err := cl.dockerClient.DoAuthRequest(req)
	if err != nil {
		return err
	}

	cl.token = resp.Token
	return nil
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (cl *Client) listTags(name string) ([]string, error) {
	req, err := cl.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cl.token))

	var tagsResp listTagsResponse
	if err := cl.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	return tagsResp.Tags, nil
}
