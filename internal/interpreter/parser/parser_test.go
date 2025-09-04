package parser

import (
	"errors"
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantCmd *Command
		wantErr error
	}{
		{
			name:    "valid SET command",
			input:   "SET weather_2_pm cold_moscow_weather",
			wantCmd: &Command{CommandType: SetCommandType, Arguments: []string{"weather_2_pm", "cold_moscow_weather"}},
			wantErr: nil,
		},
		{
			name:    "valid SET command (tabs)",
			input:   "SET\tweather_2_pm\t\tcold_moscow_weather",
			wantCmd: &Command{CommandType: SetCommandType, Arguments: []string{"weather_2_pm", "cold_moscow_weather"}},
			wantErr: nil,
		},
		{
			name:    "valid GET command",
			input:   "GET /etc/nginx/config",
			wantCmd: &Command{CommandType: GetCommandType, Arguments: []string{"/etc/nginx/config"}},
			wantErr: nil,
		},
		{
			name:    "valid DEL command",
			input:   "DEL user_****",
			wantCmd: &Command{CommandType: DelCommandType, Arguments: []string{"user_****"}},
			wantErr: nil,
		},
		{
			name:    "SET command not enough arguments",
			input:   "SET weather_2_pm",
			wantCmd: nil,
			wantErr: ErrNoEnoughArgumentsForSetCommand,
		},
		{
			name:    "GET command not enough arguments",
			input:   "GET",
			wantCmd: nil,
			wantErr: ErrNoEnoughArgumentsForGetCommand,
		},
		{
			name:    "DEL command not enough arguments",
			input:   "DEL",
			wantCmd: nil,
			wantErr: ErrNoEnoughArgumentsForDelCommand,
		},
		{
			name:    "unknown command",
			input:   "UPDATE something",
			wantCmd: nil,
			wantErr: ErrUnknownCommandType,
		},
		{
			name:    "empty query",
			input:   "",
			wantCmd: nil,
			wantErr: ErrNoTokensInQuery,
		},
		{
			name:    "spaces only",
			input:   "    ",
			wantCmd: nil,
			wantErr: ErrNoTokensInQuery,
		},
		{
			name:    "extra arguments for SET",
			input:   "SET key value extra",
			wantCmd: &Command{CommandType: SetCommandType, Arguments: []string{"key", "value"}},
			wantErr: nil,
		},
		{
			name:    "extra arguments for GET",
			input:   "GET key extra",
			wantCmd: &Command{CommandType: GetCommandType, Arguments: []string{"key"}},
			wantErr: nil,
		},
		{
			name:    "extra arguments for DEL",
			input:   "DEL key extra",
			wantCmd: &Command{CommandType: DelCommandType, Arguments: []string{"key"}},
			wantErr: nil,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotErr := parser.Parse(tt.input)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Parse() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t.Errorf("Parse() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}
