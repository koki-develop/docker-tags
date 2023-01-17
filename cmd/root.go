package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/distribution/distribution/reference"
	"github.com/koki-develop/docker-tags/pkg/registry/artifactregistry"
	"github.com/koki-develop/docker-tags/pkg/registry/dockerhub"
	"github.com/koki-develop/docker-tags/pkg/registry/ecr"
	"github.com/koki-develop/docker-tags/pkg/registry/ecrpublic"
	"github.com/koki-develop/docker-tags/pkg/registry/gcr"
	"github.com/spf13/cobra"
)

type client interface {
	ListTags(name string) ([]string, error)
}

var (
	_ client = (*dockerhub.Client)(nil)
	_ client = (*ecr.Client)(nil)
	_ client = (*ecrpublic.Client)(nil)
	_ client = (*artifactregistry.Client)(nil)
	_ client = (*gcr.Client)(nil)
)

var (
	awsProfile string
)

func newClient(domain string) (client, error) {
	switch {
	case domain == "docker.io":
		return dockerhub.New(), nil
	case domain == "public.ecr.aws":
		return ecrpublic.New(&ecrpublic.Config{
			Profile: awsProfile,
		}), nil
	case strings.HasSuffix(domain, "amazonaws.com"):
		// <AWS_ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/<REPOSITORY_NAME>
		return ecr.New(&ecr.Config{
			Profile: awsProfile,
			Domain:  domain,
		}), nil
	case domain == "gcr.io":
		return gcr.New(), nil
	case strings.HasSuffix(domain, "-docker.pkg.dev"):
		// <LOCATION>-docker.pkg.dev/<PROJECT>/<REPOSITORY>/<PACKAGE>
		return artifactregistry.New(&artifactregistry.Config{Domain: domain}), nil
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

func init() {
	rootCmd.Flags().StringVar(&awsProfile, "aws-profile", "", "aws profile")
}
