package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

type (
	Timeout  time.Duration
	DataSize uint64
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
	TB = 1 << 40
)

func (t *Timeout) UnmarshalYAML(value *yaml.Node) error {
	d, err := time.ParseDuration(value.Value)
	if err != nil {
		return fmt.Errorf("can't parse timeout: %w", err)
	}
	*t = Timeout(d)
	return nil
}

func (d *DataSize) UnmarshalYAML(value *yaml.Node) error {
	re := regexp.MustCompile(`(?i)^(\d+)([a-zA-Z]+)?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(value.Value))
	if len(matches) < 2 {
		return errors.New("invalid data size format")
	}

	v, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid data size format: %w", err)
	}

	unit := strings.ToLower(matches[2])
	switch unit {
	case "", "b":
		*d = DataSize(v)
	case "kb":
		*d = DataSize(v * KB)
	case "mb":
		*d = DataSize(v * MB)
	case "gb":
		*d = DataSize(v * GB)
	case "tb":
		*d = DataSize(v * TB)
	default:
		return fmt.Errorf("invalid data size format: unknown unit: %s", unit)
	}

	return nil
}

type ConfigEngine struct {
	Type string `yaml:"type"`
}

type ConfigNetwork struct {
	Address        string   `yaml:"address"`
	MaxConnections int      `yaml:"max_connections"`
	MaxMessageSize DataSize `yaml:"max_message_size"`
	IdleTimeout    Timeout  `yaml:"idle_timeout"`
}

type ConfigLogging struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type ConfigWAL struct {
	FlushingBatchSize    int      `yaml:"flushing_batch_size"`
	FlushingBatchTimeout Timeout  `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       DataSize `yaml:"max_segment_size"`
	DataDirectory        string   `yaml:"data_directory"`
}

type Config struct {
	Engine  ConfigEngine  `yaml:"engine"`
	Network ConfigNetwork `yaml:"network"`
	Logging ConfigLogging `yaml:"logging"`
	WAL     ConfigWAL     `yaml:"wal"`
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
			MaxMessageSize: DataSize(4 * KB),
			IdleTimeout:    Timeout(5 * time.Minute),
		},
		Logging: ConfigLogging{
			Level:  LoggingLevelInfo,
			Output: os.Stderr.Name(),
		},
		WAL: ConfigWAL{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: Timeout(10 * time.Millisecond),
			MaxSegmentSize:       DataSize(10 * MB),
			DataDirectory:        "wal",
		},
	}
}
