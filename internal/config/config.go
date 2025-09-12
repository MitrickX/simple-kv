package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Engine struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`
	Network struct {
		Address        string `yaml:"address"`
		MaxConnections int    `yaml:"max_connections"`
		MaxMessageSize string `yaml:"max_message_size"`
		IdleTimeout    string `yaml:"idle_timeout"`
	} `yaml:"network"`
	Logging struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	} `yaml:"logging"`
}

// Parse reads a YAML config file and unmarshals it into Config.
func Parse(filePath string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
