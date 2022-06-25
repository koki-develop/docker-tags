package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Tag struct {
	Name string `json:"name"`
}

type Tags []*Tag

type Client struct {
	httpAPI HTTPAPI
}

func NewClient() *Client {
	return &Client{
		httpAPI: new(http.Client),
	}
}

func (c *Client) FetchTags(img string) (Tags, error) {
	u, err := url.ParseRequestURI("https://registry.hub.docker.com/v1/repositories")
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, img, "tags")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpAPI.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// pass
	case http.StatusNotFound:
		return nil, errors.New("resource not found")
	default:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("request error: %s", string(b))
	}

	var ts Tags
	if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
		return nil, err
	}

	return ts, nil
}
