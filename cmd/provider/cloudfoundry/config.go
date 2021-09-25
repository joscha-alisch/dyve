package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database     DatabaseConfig
	Port         int      `yaml:"port"`
	CloudFoundry CfConfig `yaml:"cloudfoundry"`
}

type CfConfig struct {
	Api      string `yaml:"api"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type DatabaseConfig struct {
	URI  string `yaml:"uri"`
	Name string `yaml:"name"`
}

func LoadFrom(path string) (Config, error) {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("dyve")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	viper.AutomaticEnv()

	c := Config{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return Config{}, err
	}
	return c, err
}
