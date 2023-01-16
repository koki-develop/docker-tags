package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/distribution/distribution/reference"
	"github.com/koki-develop/docker-tags/pkg/docker"
	"github.com/koki-develop/docker-tags/pkg/ecr"
	"github.com/spf13/cobra"
)

type client interface {
	ListTags(name string) ([]string, error)
}

var (
	_ client = (*docker.Client)(nil)
	_ client = (*ecr.Client)(nil)
)

func newClient(domain string) (client, error) {
	switch {
	case domain == "docker.io":
		return docker.New(), nil
	case strings.HasSuffix(domain, "amazonaws.com"):
		return ecr.New(), nil
	default:
		return nil, fmt.Errorf("unsupported image repository: %s", domain)
	}
}

var rootCmd = &cobra.Command{
	Use:   "docker-tags [IMAGE]",
	Short: "Command line tool to get a list of tags for docker images.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		img := args[0]

		named, err := reference.ParseNormalizedNamed(img)
		if err != nil {
			return err
		}

		d := reference.Domain(named)
		p := reference.Path(named)

		cl, err := newClient(d)
		if err != nil {
			return err
		}

		tags, err := cl.ListTags(p)
		if err != nil {
			return err
		}

		for _, t := range tags {
			fmt.Println(t)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
