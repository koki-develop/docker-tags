package googleutil

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (cl *Client) Token() (*oauth2.Token, error) {
	ctx := context.Background()
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	tkn, err := cred.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return tkn, nil
}
