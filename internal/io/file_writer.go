package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func WriteLines(filename string, lines []string) error {
	var w io.Writer
	if filename == "" || filename == "-" {
		w = os.Stdout
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", filename, err)
		}
		defer f.Close()
		w = f
	}

	bw := bufio.NewWriter(w)
	for _, line := range lines {
		if _, err := bw.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return bw.Flush()
}
