package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	authURL    string
	httpClient *http.Client
}

type Config struct {
	AuthURL string
}

func New(cfg *Config) *Client {
	return &Client{
		authURL:    cfg.AuthURL,
		httpClient: new(http.Client),
	}
}

type DockerAuthResponse struct {
	Token string `json:"token"`
}

func (cl *Client) NewAuthRequest(name string) (*http.Request, error) {
	u, err := url.ParseRequestURI(cl.authURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("service", "registry.docker.io")
	q.Set("scope", fmt.Sprintf("repository:%s:pull", name))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (cl *Client) DoAuthRequest(req *http.Request) (*DockerAuthResponse, error) {
	var resp DockerAuthResponse
	if err := cl.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cl *Client) do(req *http.Request, out interface{}) error {
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

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}

	return nil
}
