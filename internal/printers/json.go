package printers

import (
	"encoding/json"
	"fmt"
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
	elms := make([]string, len(params.Tags))

	if params.WithName {
		for i, t := range params.Tags {
			elms[i] = fmt.Sprintf("%s:%s", params.Image, t)
		}
	} else {
		copy(elms, params.Tags)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(elms); err != nil {
		return err
	}

	return nil
}
