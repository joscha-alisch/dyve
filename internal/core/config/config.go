package config

import (
	"bytes"
	"github.com/fatih/structs"
	"github.com/jeremywohl/flatten"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"strings"
)

type Config struct {
	LogLevel       string           `yaml:"logLevel"`
	DevConfig      DevConfig        `yaml:"devConfig"`
	Providers      []ProviderConfig `yaml:"providers"`
	Database       DatabaseConfig   `yaml:"database"`
	Port           int              `yaml:"port"`
	Reconciliation ReconConfig      `yaml:"reconciliation"`
	Auth           AuthConfig       `yaml:"auth"`
	ExternalUrl    string           `yaml:"externalUrl"`
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

func defaults() Config {
	return Config{
		LogLevel: "info",
		DevConfig: DevConfig{
			UseFakeOauth2:      false,
			DisableAuth:        false,
			UserGroups:         nil,
			DisableOriginCheck: false,
		},
		Providers: nil,
		Database: DatabaseConfig{
			URI:  "mongodb://localhost:27017",
			Name: "dyve_core",
		},
		Port: 9000,
		Reconciliation: ReconConfig{
			CacheSeconds: 20,
		},
		Auth: AuthConfig{
			Secret: "",
			GitHub: AuthProviderConfig{
				Enabled: false,
			},
		},
		ExternalUrl: "http://localhost:9000",
	}
}

func LoadFrom(path string) (Config, error) {
	v := viper.New()
	b, err := yaml.Marshal(defaults())
	if err != nil {
		return Config{}, err
	}
	defaultConfig := bytes.NewReader(b)
	v.SetConfigType("yaml")
	if err := v.MergeConfig(defaultConfig); err != nil {
		return Config{}, err
	}

	v.SetConfigFile(path)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return Config{}, err
		}
	}

	v.AutomaticEnv()
	v.SetEnvPrefix("dyve")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	confMap := structs.Map(Config{})

	// Flatten nested conf map
	flat, err := flatten.Flatten(confMap, "", flatten.DotStyle)
	if err != nil {
		return Config{}, errors.Wrap(err, "Unable to flatten config")
	}

	// Bind each conf fields to environment vars
	for key, _ := range flat {
		err := v.BindEnv(key)
		if err != nil {
			return Config{}, errors.Wrapf(err, "Unable to bind env var: %s", key)
		}
	}

	config := &Config{}
	err = v.Unmarshal(&config)
	return *config, err
}
