package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

func showHostConfig(cfg *Config) (Host, error) {
	args := flag.Args()
	if len(args) == 0 {
		return Host{}, errors.New("Need specify connection name")
	}
	if len(args) != 1 {
		return Host{}, errors.New("Too much input arguments. Need specify only one connection name")
	}

	h, err := cfg.getHost(args[0])
	if err != nil {
		return h, err
	}

	home, err := homeDir()
	if err != nil {
		return h, err
	}

	fileName := filepath.Join(home, h.fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return h, fmt.Errorf("Can't open config file %q: %w", fileName, err)
	}

	var text []string
	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	for i := 1; s.Scan(); i++ {
		if i >= h.lineStart {
			text = append(text, s.Text())
		}
		if i > h.lineEnd {
			break
		}
	}

	h.text = strings.Join(text, "\n")

	return h, nil
}
