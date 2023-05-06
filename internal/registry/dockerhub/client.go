package dockerhub

import (
	"fmt"
	"net/http"

	"github.com/koki-develop/docker-tags/internal/util/dockerutil"
)

type Registry struct {
	dockerClient *dockerutil.Client
	httpClient   *http.Client
}

func New() *Registry {
	return &Registry{
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  "https://registry.hub.docker.com",
			AuthURL: "https://auth.docker.io/token",
		}),
		httpClient: new(http.Client),
	}
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
	req, err := r.dockerClient.NewAuthRequest(name)
	if err != nil {
		return "", err
	}

	resp, err := r.dockerClient.DoAuthRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (r *Registry) listTags(name, tkn string) ([]string, error) {
	req, err := r.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))

	var tagsResp listTagsResponse
	if err := r.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	// reverse
	tags := tagsResp.Tags
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}

	return tags, nil
}
