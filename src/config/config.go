package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	defaultStoragePath = "storage.db"
	defaultSchemaPath  = "schema.json"
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
	Path       string `yaml:"path"`
	File       string `yaml:"file"`
	Period     string `yaml:"period"`
	Build      *Build `yaml:"build"`
	Run        *Run   `yaml:"run"`
	LogsEnable bool   `yaml:"log"`
}

type Database struct {
	Enable  bool   `yaml:"enable"`
	Storage string `yaml:"storage" default:"storage.db"`
	Schema  string `yaml:"schema" default:"schema.json"`
}

type Config struct {
	Url      string                          `yaml:"url"`
	LogPath  string                          `yaml:"log"`
	LogLevel string                          `yaml:"log_level"`
	Database Database                        `yaml:"database"`
	Handlers map[string]ExecutableServerUnit `yaml:"handlers"`
	Tasks    map[string]ExecutableServerUnit `yaml:"tasks"`
}

func (config *Config) setDefaults() {
	if config.Database.Storage == "" {
		config.Database.Storage = defaultStoragePath
	}

	if config.Database.Schema == "" {
		config.Database.Schema = defaultSchemaPath
	}
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

	config.setDefaults()

	fmt.Println("[Config] config loaded successfully: ", *config)

	return config, nil
}
