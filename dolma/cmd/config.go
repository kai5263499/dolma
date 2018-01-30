package cmd

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TargetBinary 		  string `config:"required;usage=Binary to add resources to"`
	SectionPrefix	      string `config:"usage=Path prefix to use for resources"`
	SectionType           string `config:"default=zip;usage=Type of section to add (zip)"`
	SectionContent 		  string `config:"usage=File or folder to add to section prefix"`
}

func (c *Config) DumpJSON() (string, error) {
	v, err := json.Marshal(c)
	return string(v), err
}

func (c *Config) DumpYAML() (string, error) {
	v, err := yaml.Marshal(c)
	return string(v), err
}

func (c *Config) Validate() error {
	return nil
}
