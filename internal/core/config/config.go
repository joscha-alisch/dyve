package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	AppProviders []AppProviderConfig `yaml:"appProviders"`
	Database     DatabaseConfig
	Port         int `yaml:"port"`
}

type AppProviderConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

func LoadFrom(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	return Load(f)
}

func Load(r io.Reader) (Config, error) {
	c := Config{}
	err := yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
