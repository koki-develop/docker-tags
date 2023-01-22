package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type printer interface {
	Print(params *printParams) error
}

type printParams struct {
	Name     string
	Tags     []string
	WithName bool
}

func newPrinter(output string) (printer, error) {
	switch output {
	case "text":
		return &textPrinter{w: os.Stdout}, nil
	case "json":
		return &jsonPrinter{w: os.Stdout}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", output)
	}
}

var (
	_ printer = (*textPrinter)(nil)
	_ printer = (*jsonPrinter)(nil)
)

type textPrinter struct {
	w io.Writer
}

func (p *textPrinter) Print(params *printParams) error {
	if params.WithName {
		for _, t := range params.Tags {
			if _, err := fmt.Fprintf(p.w, "%s:%s\n", params.Name, t); err != nil {
				return err
			}
		}
	} else {
		for _, t := range params.Tags {
			if _, err := fmt.Fprintln(p.w, t); err != nil {
				return err
			}
		}
	}

	return nil
}

type jsonPrinter struct {
	w io.Writer
}

func (p *jsonPrinter) Print(params *printParams) error {
	if params.WithName {
		for i, t := range params.Tags {
			params.Tags[i] = fmt.Sprintf("%s:%s", params.Name, t)
		}
	}

	j, err := json.MarshalIndent(params.Tags, "", "  ")
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(p.w, string(j)); err != nil {
		return err
	}

	return nil
}
