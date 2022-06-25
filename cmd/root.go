package cmd

import (
	"os"

	"github.com/koki-develop/docker-tags/pkg/docker"
	"github.com/koki-develop/docker-tags/pkg/report"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "docker-tags IMAGE",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		img := args[0]
		ot, err := cmd.Flags().GetBool("only-tags")
		if err != nil {
			return err
		}

		c := docker.NewClient()
		tags, err := c.FetchTags(img)
		if err != nil {
			return err
		}

		r := report.New(&report.Options{
			Writer:   os.Stdout,
			OnlyTags: ot,
		})
		return r.Print(img, tags)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("only-tags", "t", false, "show tags only")
}
