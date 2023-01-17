package artifactregistry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	authURL := url.URL{Scheme: "https", Path: path.Join(cfg.Domain, "/v2/token")}
	return &Client{
		domain: cfg.Domain,
		dockerClient: docker.New(&docker.Config{
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
	u := url.URL{Scheme: "https", Host: cl.domain}
	u.Path = path.Join(u.Path, "v2", name, "tags/list")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if err := cl.auth(name); err == nil && cl.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cl.token))
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	var tagsResp listTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
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
