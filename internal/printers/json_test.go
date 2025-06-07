package printers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSONPrinter(t *testing.T) {
	printer := NewJSONPrinter()
	assert.NotNil(t, printer)
	assert.IsType(t, &JSONPrinter{}, printer)
}

func TestJSONPrinter_Print(t *testing.T) {
	tests := []struct {
		name     string
		params   *PrintParameters
		expected string
	}{
		{
			name: "basic tags without name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{"latest", "3.18", "3.17"},
				WithName: false,
			},
			expected: "[\n  \"latest\",\n  \"3.18\",\n  \"3.17\"\n]\n",
		},
		{
			name: "basic tags with name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{"latest", "3.18", "3.17"},
				WithName: true,
			},
			expected: "[\n  \"alpine:latest\",\n  \"alpine:3.18\",\n  \"alpine:3.17\"\n]\n",
		},
		{
			name: "single tag",
			params: &PrintParameters{
				Image:    "nginx",
				Tags:     []string{"latest"},
				WithName: false,
			},
			expected: "[\n  \"latest\"\n]\n",
		},
		{
			name: "empty tags",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{},
				WithName: false,
			},
			expected: "[]\n",
		},
		{
			name: "tags with special characters",
			params: &PrintParameters{
				Image:    "myapp",
				Tags:     []string{"v1.0.0", "v1.0.0-beta", "feature-branch"},
				WithName: false,
			},
			expected: "[\n  \"v1.0.0\",\n  \"v1.0.0-beta\",\n  \"feature-branch\"\n]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := NewJSONPrinter()
			var buf bytes.Buffer

			err := printer.Print(&buf, tt.params)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}
