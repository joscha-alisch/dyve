package config

import (
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel       string           `yaml:"logLevel"`
	DevConfig      DevConfig        `yaml:"devConfig"`
	Providers      []ProviderConfig `yaml:"providers"`
	Database       DatabaseConfig
	Port           int         `yaml:"port"`
	Reconciliation ReconConfig `yaml:"reconciliation"`
	Auth           AuthConfig  `yaml:"auth"`
	ExternalUrl    string      `yaml:"externalUrl"`
}

type DevConfig struct {
	UseFakeOauth2      bool     `yaml:"useFakeOauth2"`
	DisableAuth        bool     `yaml:"disableAuth"`
	UserGroups         []string `yaml:"userGroups"`
	DisableOriginCheck bool     `yaml:"disableOriginCheck"`
}

type ProviderConfig struct {
	Id       string          `yaml:"id"`
	Name     string          `yaml:"name"`
	Host     string          `yaml:"host"`
	Features []provider.Type `yaml:"features"`
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
