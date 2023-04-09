package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func Exists(path string) (bool, error) {
	_, fileErr := os.Stat(path)
	if fileErr == nil {
		return true, nil
	}
	if os.IsNotExist(fileErr) {
		return false, nil
	}
	return true, nil
}

func AssureExists(filePath string) error {
	path := filepath.Dir(filePath)
	exists, err := Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Couldn't create path: %s", path)
		}
	}
	return nil
}
