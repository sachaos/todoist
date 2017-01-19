package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type Writer interface {
	Write([]string) error
	Flush()
}

type TSVWriter struct {
	w *tabwriter.Writer
}

func NewTSVWriter(w io.Writer) *TSVWriter {
	return &TSVWriter{
		w: tabwriter.NewWriter(w, 0, 4, 1, ' ', 0),
	}
}

func (w *TSVWriter) Flush() {
	w.w.Flush()
}

func (w *TSVWriter) Write(record []string) error {
	string := strings.Join(record[:], "\t")
	fmt.Fprintln(w.w, string)
	return nil
}
