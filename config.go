package main

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Host struct {
	Name     string `yaml:"-"`
	Address  string `yaml:"Address"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`

	fileName  string `yaml:"-"`
	lineStart int    `yaml:"-"`
	lineEnd   int    `yaml:"-"`
	text      string `yaml:"-"`
}

type Config struct {
	Hosts []Host
}

func (cfg *Config) UnmarshalYAML(doc *yaml.Node) error {
	for i := 0; i < len(doc.Content); i++ {
		c := doc.Content[i]
		if c.Kind != yaml.ScalarNode {
			return fmt.Errorf("Incorrect file format: line=%d, %q", c.Line, c.Value)
		}

		h := Host{Name: c.Value}
		h.lineStart = c.Line
		if c.HeadComment != "" {
			h.lineStart -= 1
		}

		i++

		m := doc.Content[i]
		if m.Kind != yaml.MappingNode {
			return fmt.Errorf("Incorrect file format: line=%d, %q", c.Line, c.Value)
		}

		if err := m.Decode(&h); err != nil {
			return err
		}

		for _, v := range m.Content {
			if h.lineEnd < v.Line {
				h.lineEnd = v.Line
			}
		}

		cfg.Hosts = append(cfg.Hosts, h)
	}
	return nil
}

func (cfg *Config) ListHosts() []string {
	names := make([]string, len(cfg.Hosts), len(cfg.Hosts))
	for i := range cfg.Hosts {
		names[i] = cfg.Hosts[i].Name
	}

	return names
}

func (cfg *Config) Parse(r io.Reader, fileName string) error {
	d := yaml.NewDecoder(r)
	d.KnownFields(true)

	var x Config
	if err := d.Decode(&x); err != nil {
		return err
	}

	for i := range x.Hosts {
		x.Hosts[i].fileName = fileName
	}

	cfg.Hosts = append(cfg.Hosts, x.Hosts...)

	return nil
}

func (cfg *Config) getHost(name string) (Host, error) {
	for i := range cfg.Hosts {
		if cfg.Hosts[i].Name == name {
			return cfg.Hosts[i], nil
		}
	}

	return Host{}, fmt.Errorf("No config for host %q", name)
}
