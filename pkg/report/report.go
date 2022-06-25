package report

import (
	"fmt"
	"io"

	"github.com/koki-develop/docker-tags/pkg/docker"
)

type Report struct {
	options *Options
}

type Options struct {
	Writer   io.Writer
	OnlyTags bool
}

func New(opts *Options) *Report {
	return &Report{
		options: opts,
	}
}

func (r *Report) Print(img string, tags docker.Tags) error {
	switch {
	case r.options.OnlyTags:
		return r.printOnlyTags(tags)
	default:
		return r.print(img, tags)
	}
}

func (r *Report) print(img string, tags docker.Tags) error {
	for _, t := range tags {
		if _, err := fmt.Fprintf(r.options.Writer, "%s:%s\n", img, t.Name); err != nil {
			return err
		}
	}
	return nil
}

func (r *Report) printOnlyTags(tags docker.Tags) error {
	for _, t := range tags {
		if _, err := fmt.Fprintln(r.options.Writer, t.Name); err != nil {
			return err
		}
	}
	return nil
}
