package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	EngineTypeInMemory  = "in_memory"
	LoggingLevelDebug   = "debug"
	LoggingLevelInfo    = "info"
	LoggingLevelWarning = "warning"
	LoggingLevelError   = "error"
	LoggingLevelPanic   = "panic"
	LoggingLevelFatal   = "fatal"
)

type Timeout time.Duration

func (t *Timeout) UnmarshalYAML(value *yaml.Node) error {
	d, err := time.ParseDuration(value.Value)
	if err != nil {
		return fmt.Errorf("can't parse timeout: %w", err)
	}
	*t = Timeout(d)
	return nil
}

type ConfigEngine struct {
	Type string `yaml:"type"`
}

type ConfigNetwork struct {
	Address        string  `yaml:"address"`
	MaxConnections int     `yaml:"max_connections"`
	MaxMessageSize string  `yaml:"max_message_size"`
	IdleTimeout    Timeout `yaml:"idle_timeout"`
}

type ConfigLogging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type Config struct {
	Engine  ConfigEngine  `yaml:"engine"`
	Network ConfigNetwork `yaml:"network"`
	Logging ConfigLogging `yaml:"logging"`
}

// Parse reads a YAML config file and unmarshals it into Config.
func Parse(filePath string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

func Default() Config {
	return Config{
		Engine: ConfigEngine{
			Type: EngineTypeInMemory,
		},
		Network: ConfigNetwork{
			Address:        "127.0.0.1:0",
			MaxConnections: 5,
			MaxMessageSize: "4KB",
			IdleTimeout:    Timeout(5 * time.Minute),
		},
		Logging: ConfigLogging{
			Level:  LoggingLevelInfo,
			Output: os.Stderr.Name(),
		},
	}
}
