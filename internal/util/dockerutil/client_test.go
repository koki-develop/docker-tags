package dockerutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config := &Config{
		APIURL:  "https://registry.example.com",
		AuthURL: "https://auth.example.com",
	}

	client := New(config)

	assert.NotNil(t, client)
	assert.Equal(t, config.APIURL, client.apiURL)
	assert.Equal(t, config.AuthURL, client.authURL)
	assert.NotNil(t, client.httpClient)
}

func TestClient_NewAuthRequest(t *testing.T) {
	client := &Client{
		authURL: "https://auth.docker.io/token",
	}

	req, err := client.NewAuthRequest("library/alpine")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "https://auth.docker.io/token?scope=repository%3Alibrary%2Falpine%3Apull&service=registry.docker.io", req.URL.String())
}

func TestClient_NewListTagsRequest(t *testing.T) {
	client := &Client{
		apiURL: "https://registry.hub.docker.com",
	}

	req, err := client.NewListTagsRequest("library/alpine")

	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.Equal(t, "https://registry.hub.docker.com/v2/library/alpine/tags/list", req.URL.String())
}
