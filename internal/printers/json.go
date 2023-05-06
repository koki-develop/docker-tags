package printers

import (
	"encoding/json"
	"io"
)

var (
	JSON = Register("json", NewJSONPrinter())
)

type JSONPrinter struct{}

func NewJSONPrinter() *JSONPrinter {
	return &JSONPrinter{}
}

func (p *JSONPrinter) Print(w io.Writer, params *PrintParameters) error {
	tags := formatTags(params)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tags); err != nil {
		return err
	}

	return nil
}
