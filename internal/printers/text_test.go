package printers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTextPrinter(t *testing.T) {
	printer := NewTextPrinter()
	assert.NotNil(t, printer)
	assert.IsType(t, &TextPrinter{}, printer)
}

func TestTextPrinter_Print(t *testing.T) {
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
			expected: "latest\n3.18\n3.17\n",
		},
		{
			name: "basic tags with name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{"latest", "3.18", "3.17"},
				WithName: true,
			},
			expected: "alpine:latest\nalpine:3.18\nalpine:3.17\n",
		},
		{
			name: "single tag",
			params: &PrintParameters{
				Image:    "nginx",
				Tags:     []string{"latest"},
				WithName: false,
			},
			expected: "latest\n",
		},
		{
			name: "empty tags",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{},
				WithName: false,
			},
			expected: "\n",
		},
		{
			name: "tags with special characters",
			params: &PrintParameters{
				Image:    "myapp",
				Tags:     []string{"v1.0.0", "v1.0.0-beta", "feature-branch"},
				WithName: false,
			},
			expected: "v1.0.0\nv1.0.0-beta\nfeature-branch\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := NewTextPrinter()
			var buf bytes.Buffer

			err := printer.Print(&buf, tt.params)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}
