package printers

import (
	"fmt"
	"io"
	"strings"
)

var (
	Text = Register("text", NewTextPrinter())
)

type TextPrinter struct{}

func NewTextPrinter() *TextPrinter {
	return &TextPrinter{}
}

func (p *TextPrinter) Print(w io.Writer, params *PrintParameters) error {
	tags := formatTags(params)
	if _, err := fmt.Fprintln(w, strings.Join(tags, "\n")); err != nil {
		return err
	}
	return nil
}
