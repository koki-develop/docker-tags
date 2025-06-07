package printers

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	originalPrinters := maps.Clone(printers)
	defer func() {
		printers = originalPrinters
	}()

	mockPrinter := &TextPrinter{}
	result := Register("mock", mockPrinter)

	assert.Equal(t, mockPrinter, result)
	assert.Equal(t, mockPrinter, printers["mock"])
}

func Test_Get(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{
			name:        "valid text format",
			format:      "text",
			expectError: false,
		},
		{
			name:        "valid json format",
			format:      "json",
			expectError: false,
		},
		{
			name:        "valid yaml format",
			format:      "yaml",
			expectError: false,
		},
		{
			name:        "invalid format",
			format:      "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := Get(tt.format)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, printer)
				assert.Contains(t, err.Error(), "unsupported output format")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, printer)
			}
		})
	}
}

func Test_List(t *testing.T) {
	formats := List()

	assert.Contains(t, formats, "text")
	assert.Contains(t, formats, "json")
	assert.Contains(t, formats, "yaml")
	assert.Len(t, formats, 3)
}

func Test_formatTags(t *testing.T) {
	tests := []struct {
		name     string
		params   *PrintParameters
		expected []string
	}{
		{
			name: "without name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{"latest", "3.18", "3.17"},
				WithName: false,
			},
			expected: []string{"latest", "3.18", "3.17"},
		},
		{
			name: "with name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{"latest", "3.18", "3.17"},
				WithName: true,
			},
			expected: []string{"alpine:latest", "alpine:3.18", "alpine:3.17"},
		},
		{
			name: "empty tags",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{},
				WithName: false,
			},
			expected: []string{},
		},
		{
			name: "empty tags with name",
			params: &PrintParameters{
				Image:    "alpine",
				Tags:     []string{},
				WithName: true,
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTags(tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}
