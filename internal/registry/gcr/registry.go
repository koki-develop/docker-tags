package gcr

import (
	"fmt"
	"sort"

	"github.com/koki-develop/docker-tags/internal/util/dockerutil"
	"github.com/koki-develop/docker-tags/internal/util/googleutil"
)

type Registry struct {
	dockerClient *dockerutil.Client
	googleClient *googleutil.Client
}

func New() *Registry {
	return &Registry{
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  "https://gcr.io",
			AuthURL: "https://gcr.io/v2/token",
		}),
		googleClient: googleutil.New(),
	}
}

type listTagsResponse struct {
	Manifest map[string]manifest `json:"manifest"`
}

type manifest struct {
	Tag            []string `json:"tag"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
}

func (r *Registry) ListTags(name string) ([]string, error) {
	token, _ := r.auth(name)

	tags, err := r.listTags(name, token)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *Registry) auth(name string) (string, error) {
	tkn, err := r.googleClient.Token()
	if err != nil {
		return "", err
	}

	req, err := r.dockerClient.NewAuthRequest(name)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("_token", tkn.AccessToken)

	resp, err := r.dockerClient.DoAuthRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (r *Registry) listTags(name, token string) ([]string, error) {
	req, err := r.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	var tagsResp listTagsResponse
	if err := r.dockerClient.DoListTagsRequest(req, &tagsResp); err != nil {
		return nil, err
	}

	manifests := []manifest{}
	for _, m := range tagsResp.Manifest {
		manifests = append(manifests, m)
	}
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].TimeUploadedMs > manifests[j].TimeUploadedMs
	})

	tags := []string{}
	for _, m := range manifests {
		tags = append(tags, m.Tag...)
	}

	return tags, nil
}
