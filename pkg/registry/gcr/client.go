package gcr

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/koki-develop/docker-tags/pkg/util/dockerutil"
	"github.com/koki-develop/docker-tags/pkg/util/google"
)

type Client struct {
	dockerClient *dockerutil.Client
	googleClient *google.Client
	httpClient   *http.Client
}

func New() *Client {
	return &Client{
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  "https://gcr.io",
			AuthURL: "https://gcr.io/v2/token",
		}),
		googleClient: google.New(),
		httpClient:   new(http.Client),
	}
}

type listTagsResponse struct {
	Manifest map[string]manifest `json:"manifest"`
}

type manifest struct {
	Tag            []string `json:"tag"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
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
	tkn, err := cl.googleClient.Token()
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
