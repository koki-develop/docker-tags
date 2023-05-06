package printers

import (
	"io"

	"gopkg.in/yaml.v3"
)

var (
	YAML = Register("yaml", NewYAMLPrinter())
)

type YAMLPrinter struct{}

func NewYAMLPrinter() *YAMLPrinter {
	return &YAMLPrinter{}
}

func (p *YAMLPrinter) Print(w io.Writer, params *PrintParameters) error {
	tags := formatTags(params)

	enc := yaml.NewEncoder(w)
	if err := enc.Encode(tags); err != nil {
		return err
	}

	return nil
}
