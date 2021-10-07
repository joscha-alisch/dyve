package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel       string              `yaml:"logLevel"`
	DevMode        bool                `yaml:"devMode"`
	AppProviders   []AppProviderConfig `yaml:"appProviders"`
	Database       DatabaseConfig
	Port           int         `yaml:"port"`
	Reconciliation ReconConfig `yaml:"reconciliation"`
	Auth           AuthConfig  `yaml:"auth"`
	ExternalUrl    string      `yaml:"externalUrl"`
}

type AppProviderConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	URI  string `yaml:"uri"`
	Name string `yaml:"name"`
}

type ReconConfig struct {
	CacheSeconds int `yaml:"cacheSeconds"`
}

type AuthConfig struct {
	Secret string             `yaml:"secret"`
	GitHub AuthProviderConfig `yaml:"github"`
}
type AuthProviderConfig struct {
	Enabled bool
	Id      string
	Secret  string
	Org     string
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
