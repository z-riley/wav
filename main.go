package main

import (
	"encoding/json"
	"errors"
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
		exitWithErr(errors.New("usage: wav <file> [flags (see wav --help)]"))
	}

	path := pflag.Arg(0)
	absPath, err := filepath.Abs(path)
	if err != nil {
		exitWithErr(fmt.Errorf("error resolving path: %w", err))
	}

	header, err := NewHeaderFromPath(absPath)
	if errors.Is(err, ErrIncorrectFormat) {
		exitWithErr(ErrIncorrectFormat)
	} else if err != nil {
		exitWithErr(err)
	}

	switch *format {
	case FormatText:
		fmt.Println(header)

	case FormatJSON:
		j, err := json.Marshal(header)
		if err != nil {
			exitWithErr(err)
		}
		fmt.Println(string(j))

	default:
		pflag.Usage()
		os.Exit(1)
	}
}

func exitWithErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
