package artifactregistry

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	domain     string
	httpClient *http.Client
}

type Config struct {
	Domain string
}

func New(cfg *Config) *Client {
	return &Client{
		domain:     cfg.Domain,
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
