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

	"golang.org/x/oauth2/google"
)

type Client struct {
	domain     string
	httpClient *http.Client
	token      string
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

type dockerAuthResponse struct {
	Token string `json:"token"`
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

	u := url.URL{Scheme: "https", Path: path.Join(cl.domain, "/v2/token")}
	q := u.Query()
	q.Set("scope", fmt.Sprintf("repository:%s:pull", name))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth("_token", tkn.AccessToken)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}

	var authResp dockerAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	cl.token = authResp.Token
	return nil
}
