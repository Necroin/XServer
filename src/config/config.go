package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Build struct {
	Tool  string   `yaml:"tool"`
	Flags []string `yaml:"flags"`
}

type Run struct {
	Tool string   `yaml:"tool"`
	Args []string `yaml:"arguments"`
}

type ExecutableServerUnit struct {
	Path   string `yaml:"path"`
	File   string `yaml:"file"`
	Period string `yaml:"period"`
	Build  *Build `yaml:"build"`
	Run    *Run   `yaml:"run"`
}

type Config struct {
	Url      string                          `yaml:"url"`
	LogsPath string                          `yaml:"log"`
	Handlers map[string]ExecutableServerUnit `yaml:"handlers"`
	Tasks    map[string]ExecutableServerUnit `yaml:"tasks"`
}

func Load(path string) (*Config, error) {
	fmt.Printf("[Config] read config file: %s\n", path)

	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[Config] [Error] failed read config file: %s\n", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("[Config] [Error] failed map config file: %s\n", err)
	}

	fmt.Println("[Config] config loaded successfully: ", config)
	return config, nil
}
