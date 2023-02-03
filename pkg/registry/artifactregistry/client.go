package artifactregistry

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"

	"github.com/koki-develop/docker-tags/pkg/util/dockerutil"
	"github.com/koki-develop/docker-tags/pkg/util/googleutil"
)

type Client struct {
	dockerClient *dockerutil.Client
	googleClient *googleutil.Client
	httpClient   *http.Client
}

type Config struct {
	Domain string
}

func New(cfg *Config) *Client {
	apiURL := url.URL{Scheme: "https", Path: cfg.Domain}
	authURL := url.URL{Scheme: "https", Path: path.Join(cfg.Domain, "/v2/token")}
	return &Client{
		dockerClient: dockerutil.New(&dockerutil.Config{
			APIURL:  apiURL.String(),
			AuthURL: authURL.String(),
		}),
		httpClient: new(http.Client),
	}
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

type listTagsResponse struct {
	Manifest map[string]manifest `json:"manifest"`
}

type manifest struct {
	Tag            []string `json:"tag"`
	TimeUploadedMs string   `json:"timeUploadedMs"`
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
