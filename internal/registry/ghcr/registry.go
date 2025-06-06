package ghcr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Registry struct {
	httpClient *http.Client
}

func New() *Registry {
	return &Registry{
		httpClient: new(http.Client),
	}
}

type tokenResponse struct {
	Token string `json:"token"`
}

type listTagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
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
	u, err := url.ParseRequestURI("https://ghcr.io/token")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("scope", fmt.Sprintf("repository:%s:pull", name))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(b))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}

func (r *Registry) listTags(name, tkn string) ([]string, error) {
	u, err := url.ParseRequestURI("https://ghcr.io")
	if err != nil {
		return nil, err
	}
	u.Path = fmt.Sprintf("/v2/%s/tags/list", name)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tkn))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

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

	tags := tagsResp.Tags
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}

	return tags, nil
}
