package cmd

import (
	"fmt"
	"os"

	"github.com/distribution/distribution/reference"
	"github.com/koki-develop/docker-tags/pkg/docker"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "docker-tags",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		img := args[0]

		named, err := reference.ParseNormalizedNamed(img)
		if err != nil {
			return nil
		}

		domain := reference.Domain(named)
		name := reference.Path(named)

		switch domain {
		case "docker.io":
			cl := docker.New()
			tags, err := cl.ListTags(name)
			if err != nil {
				return err
			}
			for _, t := range tags {
				fmt.Println(t)
			}
		default:
			return fmt.Errorf("unsupported image repository: %s", domain)
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
