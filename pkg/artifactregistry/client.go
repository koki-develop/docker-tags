package artifactregistry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/koki-develop/docker-tags/pkg/docker"
	"golang.org/x/oauth2/google"
)

type Client struct {
	domain string
	token  string

	dockerClient *docker.Client
	httpClient   *http.Client
}

type Config struct {
	Domain string
}

func New(cfg *Config) *Client {
	apiURL := url.URL{Scheme: "https", Path: cfg.Domain}
	authURL := url.URL{Scheme: "https", Path: path.Join(cfg.Domain, "/v2/token")}
	return &Client{
		domain: cfg.Domain,
		dockerClient: docker.New(&docker.Config{
			APIURL:  apiURL.String(),
			AuthURL: authURL.String(),
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
	req, err := cl.dockerClient.NewListTagsRequest(name)
	if err != nil {
		return nil, err
	}
	if err := cl.auth(name); err == nil && cl.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cl.token))
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

func (cl *Client) auth(name string) error {
	ctx := context.Background()
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return err
	}

	tkn, err := cred.TokenSource.Token()
	if err != nil {
		return err
	}

	req, err := cl.dockerClient.NewAuthRequest(name)
	if err != nil {
		return err
	}
	req.SetBasicAuth("_token", tkn.AccessToken)

	resp, err := cl.dockerClient.DoAuthRequest(req)
	if err != nil {
		return err
	}

	cl.token = resp.Token
	return nil
}
