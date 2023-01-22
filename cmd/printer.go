package cmd

import (
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
	default:
		return nil, fmt.Errorf("unsupported output format: %s", output)
	}
}

var (
	_ printer = (*textPrinter)(nil)
)

type textPrinter struct {
	w io.Writer
}

func (p *textPrinter) Print(params *printParams) error {
	if params.WithName {
		for _, t := range params.Tags {
			fmt.Fprintf(p.w, "%s:%s\n", params.Name, t)
		}
	} else {
		for _, t := range params.Tags {
			fmt.Fprintln(p.w, t)
		}
	}

	return nil
}
