package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewHeader(t *testing.T) {
	entries, err := os.ReadDir("./samples/")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		path := "./samples/" + entry.Name()
		if !strings.Contains(path, ".wav") {
			continue
		}

		t.Run(path, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, 44)
			_, err = io.ReadAtLeast(f, buf, 44)
			if err != nil {
				t.Fatal(err)
			}

			h, err := NewHeader(buf)
			if err != nil {
				t.Error(err)
			}

			if h.Size == 0 {
				t.Error("expected size greater than 0")
			}
		})
	}
}
