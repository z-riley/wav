package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

const (
	FormatText string = "text"
	FormatJSON string = "json"
)

func main() {
	format := pflag.StringP("format", "f", FormatText, "output format (text|json)")
	pflag.Parse()

	if pflag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: wav <file> [flags (see wav --help)]")
		os.Exit(1)
	}

	path := pflag.Arg(0)
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving path:", err)
		os.Exit(1)
	}

	header, err := NewHeaderFromPath(absPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch *format {
	case FormatText:
		fmt.Println(header)

	case FormatJSON:
		j, err := json.Marshal(header)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(string(j))

	default:
		pflag.Usage()
		os.Exit(1)
	}
}
