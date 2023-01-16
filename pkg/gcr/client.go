package gcr

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: new(http.Client),
	}
}

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (cl *Client) ListTags(name string) ([]string, error) {
	u, err := url.ParseRequestURI("https://gcr.io/v2")
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, name, "tags/list")

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

	return tagsResp.Tags, nil
}
