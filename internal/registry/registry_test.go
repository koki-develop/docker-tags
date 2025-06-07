package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koki-develop/docker-tags/internal/registry/artifactregistry"
	"github.com/koki-develop/docker-tags/internal/registry/dockerhub"
	"github.com/koki-develop/docker-tags/internal/registry/ecr"
	"github.com/koki-develop/docker-tags/internal/registry/ecrpublic"
	"github.com/koki-develop/docker-tags/internal/registry/gcr"
	"github.com/koki-develop/docker-tags/internal/registry/ghcr"
)

func Test_New(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		config        *Config
		expectedType  any
		expectError   bool
		errorContains string
	}{
		{
			name:         "docker hub",
			domain:       "docker.io",
			config:       &Config{},
			expectedType: &dockerhub.Registry{},
		},
		{
			name:         "ecr public",
			domain:       "public.ecr.aws",
			config:       &Config{AWSProfile: "default"},
			expectedType: &ecrpublic.Registry{},
		},
		{
			name:         "private ecr",
			domain:       "123456789012.dkr.ecr.us-west-2.amazonaws.com",
			config:       &Config{AWSProfile: "default"},
			expectedType: &ecr.Registry{},
		},
		{
			name:         "private ecr different region",
			domain:       "987654321098.dkr.ecr.eu-west-1.amazonaws.com",
			config:       &Config{AWSProfile: "test"},
			expectedType: &ecr.Registry{},
		},
		{
			name:         "google container registry",
			domain:       "gcr.io",
			config:       &Config{},
			expectedType: &gcr.Registry{},
		},
		{
			name:         "github container registry",
			domain:       "ghcr.io",
			config:       &Config{},
			expectedType: &ghcr.Registry{},
		},
		{
			name:         "google artifact registry",
			domain:       "us-central1-docker.pkg.dev",
			config:       &Config{},
			expectedType: &artifactregistry.Registry{},
		},
		{
			name:         "google artifact registry different region",
			domain:       "europe-west1-docker.pkg.dev",
			config:       &Config{},
			expectedType: &artifactregistry.Registry{},
		},
		{
			name:          "unsupported registry",
			domain:        "registry.example.com",
			config:        &Config{},
			expectError:   true,
			errorContains: "unsupported image registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry, err := New(tt.domain, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, registry)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, registry)
				assert.IsType(t, tt.expectedType, registry)
				// Also verify it implements the Registry interface
				assert.Implements(t, (*Registry)(nil), registry)
			}
		})
	}
}
