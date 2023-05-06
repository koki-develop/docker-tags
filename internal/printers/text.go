package printers

import (
	"fmt"
	"io"
)

var (
	Text = Register("text", NewTextPrinter())
)

type TextPrinter struct{}

func NewTextPrinter() *TextPrinter {
	return &TextPrinter{}
}

func (p *TextPrinter) Print(w io.Writer, params *PrintParameters) error {
	if params.WithName {
		for _, t := range params.Tags {
			if _, err := fmt.Fprintf(w, "%s:%s\n", params.Image, t); err != nil {
				return err
			}
		}
	} else {
		for _, t := range params.Tags {
			if _, err := fmt.Fprintln(w, t); err != nil {
				return err
			}
		}
	}

	return nil
}
