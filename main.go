package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

func usage() {
	w := os.Stderr
	fmt.Fprintln(w, "dssh")
	fmt.Fprintf(w, "Simple console SSH connection manager with password authentication\n\n")
	fmt.Fprintln(w, "Usage: dssh <host>                connect to host")
	fmt.Fprintln(w, "       dssh <host> --show         show host config")
	fmt.Fprintln(w, "       dssh --install-completion  inslall completion for Clink (https://github.com/chrisant996/clink)")
}

func main() {
	installCompletionFlag := flag.Bool("install-completion", false, "install completion for Clink")
	showConfigFlag := flag.Bool("show", false, "show config for given connection")

	// hidden flag
	listHostsFlag := flag.Bool("list-hosts", false, "show all hosts name")

	flag.Usage = usage
	flag.Parse()

	// get config directory name
	home, err := homeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// install CLink completions
	if *installCompletionFlag {
		if err := installCompletion(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "See more details:", "https://github.com/chrisant996/clink")
			os.Exit(1)
		}
		fmt.Println("Need restart shell session to reload completion")
		return
	}

	cfg, err := parseConfigFiles(home)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Provide list of hosts for CLink completions
	// This is hidden command
	if *listHostsFlag {
		for _, h := range cfg.ListHosts() {
			fmt.Println(h)
		}
		return
	}

	// Show host config
	if *showConfigFlag {
		host, err := showHostConfig(cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("File:", host.fileName)
		fmt.Println(host.text)
		return

	}

	host, err := cfg.getHost(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	addr, username, pass := host.Address, host.UserName, host.Password

	if err := connect(host.Name, addr, username, pass); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

const dsshHomeDir = ".dssh"

func homeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("Can't find user's home directory: %w", err)
		return "", err
	}

	return filepath.Join(home, dsshHomeDir), nil
}

func parseConfigFiles(dir string) (*Config, error) {
	var cfg Config

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Can't open config directory (%s): %w", dir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		if filepath.Ext(e.Name()) != ".hosts" {
			continue
		}

		fileName := filepath.Join(dir, e.Name())
		file, err := os.Open(fileName)
		if err != nil {
			return nil, fmt.Errorf("Can't open hosts configuration file %q: %w", fileName, err)
		}

		if err := cfg.Parse(file, e.Name()); err != nil {
			err = fmt.Errorf("Config file %q: %s", e.Name(), err)
			return nil, err
		}
	}

	return &cfg, nil
}
