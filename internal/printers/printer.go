package printers

import (
	"fmt"
	"io"
)

type PrintParameters struct {
	Image    string
	Tags     []string
	WithName bool
}

type Printer interface {
	Print(w io.Writer, params *PrintParameters) error
}

var printers = map[string]Printer{}

func Register(name string, printer Printer) Printer {
	printers[name] = printer
	return printer
}

func Get(name string) (Printer, error) {
	p, ok := printers[name]
	if !ok {
		return nil, fmt.Errorf("unsupported output format: %s", name)
	}
	return p, nil
}

func List() []string {
	var names []string
	for name := range printers {
		names = append(names, name)
	}
	return names
}

func formatTags(params *PrintParameters) []string {
	if !params.WithName {
		return params.Tags
	}

	formatted := make([]string, len(params.Tags))
	for i, tag := range params.Tags {
		formatted[i] = fmt.Sprintf("%s:%s", params.Image, tag)
	}
	return formatted
}
