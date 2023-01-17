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

type Client struct {
	token      string
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: new(http.Client),
	}
}

func (cl *Client) ListTags(name string) ([]string, error) {
	if err := cl.auth(name); err != nil {
		return nil, err
	}

	tags, err := cl.listTags(name)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

type dockerAuthResponse struct {
	Token string `json:"token"`
}

func (cl *Client) auth(name string) error {
	u, err := url.ParseRequestURI("https://auth.docker.io/token")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("service", "registry.docker.io")
	q.Set("scope", fmt.Sprintf("repository:%s:pull", name))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

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

type listTagsResponse struct {
	Tags []string `json:"tags"`
}

func (cl *Client) listTags(name string) ([]string, error) {
	u, err := url.ParseRequestURI("https://registry.hub.docker.com/v2/")
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, name, "tags/list")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cl.token))

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
