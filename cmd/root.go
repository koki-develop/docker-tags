package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
)

type Tag struct {
	Name string `json:"name"`
}

type Tags []*Tag

var rootCmd = &cobra.Command{
	Use:  "docker-tags IMAGE",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		img := args[0]
		ot, err := cmd.Flags().GetBool("only-tags")
		if err != nil {
			return err
		}

		u, err := url.ParseRequestURI("https://registry.hub.docker.com/v1/repositories")
		if err != nil {
			return err
		}
		u.Path = path.Join(u.Path, img, "tags")

		c := new(http.Client)
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return err
		}
		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
		case http.StatusNotFound:
			return errors.New("Resource not found")
		default:
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("Request error: %s", string(b))
		}

		var ts Tags
		if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
			return err
		}

		if ot {
			for _, t := range ts {
				fmt.Println(t.Name)
			}
		} else {
			for _, t := range ts {
				fmt.Printf("%s:%s\n", img, t.Name)
			}
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
	rootCmd.Flags().BoolP("only-tags", "t", false, "show tags only")
}
