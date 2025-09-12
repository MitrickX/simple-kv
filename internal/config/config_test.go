package config

import (
    "os"
    "testing"
)

func TestParse(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        wantErr  bool
        wantType string
    }{
        {
            name: "valid config",
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
            wantErr:  false,
            wantType: "in_memory",
        },
        {
            name: "missing engine type",
            content: `
engine: {}
network:
  address: "127.0.0.1:0"
  max_connections: 100
  max_message_size: "4KB"
  idle_timeout: 5m
logging:
  level: "info"
  output: "/log/output.log"
`,
            wantErr:  false,
            wantType: "",
        },
        {
            name: "invalid yaml",
            content: `
engine
  type: "in_memory"
`,
            wantErr:  true,
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
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if cfg.Engine.Type != tt.wantType {
                    t.Errorf("expected engine type %q, got %q", tt.wantType, cfg.Engine.Type)
                }
            }
        })
    }
}