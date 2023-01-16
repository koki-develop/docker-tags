package cmd

import (
	"fmt"
	"os"

	"github.com/distribution/distribution/reference"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "docker-tags",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		img := args[0]

		name, err := reference.ParseNormalizedNamed(img)
		if err != nil {
			return nil
		}

		fmt.Printf("%#v\n", name)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
