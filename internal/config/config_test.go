package config

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantConfig Config
		wantErr    error
	}{
		{
			name: "valid_config",
			content: `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:0"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`,
			wantConfig: Config{
				Engine: ConfigEngine{
					Type: "in_memory",
				},
				Network: ConfigNetwork{
					Address:        "127.0.0.1:0",
					MaxConnections: 100,
					MaxMessageSize: DataSize(4 * KB),
					IdleTimeout:    Timeout(5 * time.Minute),
				},
				Logging: ConfigLogging{
					Level:  "info",
					Output: "/log/output.log",
				},
			},
		},
		{
			name: "invalid_timeout",
			content: `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:0"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: "invalid_timeout"
logging:
  level: "info"
  output: "/log/output.log"
`,
			wantConfig: Config{},
			wantErr:    errors.New("can't parse timeout: time: invalid duration \"invalid_timeout\""),
		},
		{
			name: "invalid_datasize",
			content: `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:0"
  max_connections: 100
  max_message_size: "invalid_size"
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`,
			wantConfig: Config{},
			wantErr:    errors.New("invalid data size format"),
		},
		{
			name: "invalid_datasize_unknown_unit",
			content: `
engine:
  type: "in_memory"
network:
  address: "127.0.0.1:0"
  max_connections: 100
  max_message_size: 100KKK
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`,
			wantConfig: Config{},
			wantErr:    errors.New("invalid data size format: unknown unit: KKK"),
		},
		{
			name: "invalid_yaml",
			content: `
engine
  type: "in_memory"
`,
			wantErr: errors.New("yaml: line 3: mapping values are not allowed in this context"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "config-*.yaml")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			cfg, err := Parse(tmpFile.Name())
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if tt.wantConfig != cfg {
					t.Errorf("expected config %+v, got %+v", tt.wantConfig, cfg)
				}
			}
		})
	}
}
