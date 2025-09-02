package interpreter

import (
	"errors"
	"reflect"
	"testing"

	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
)

func TestInterpreter_Interpret(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		mockCmd   *parser.Command
		mockErr   error
		wantRes   *Result
		wantError bool
	}{
		{
			name:      "valid SET command",
			query:     "SET weather_2_pm cold_moscow_weather",
			mockCmd:   &parser.Command{CommandType: parser.SetCommandType, Arguments: []string{"weather_2_pm", "cold_moscow_weather"}},
			mockErr:   nil,
			wantRes:   &Result{Command: parser.Command{CommandType: parser.SetCommandType, Arguments: []string{"weather_2_pm", "cold_moscow_weather"}}},
			wantError: false,
		},
		{
			name:      "valid GET command",
			query:     "GET /etc/nginx/config",
			mockCmd:   &parser.Command{CommandType: parser.GetCommandType, Arguments: []string{"/etc/nginx/config"}},
			mockErr:   nil,
			wantRes:   &Result{Command: parser.Command{CommandType: parser.GetCommandType, Arguments: []string{"/etc/nginx/config"}}},
			wantError: false,
		},
		{
			name:      "valid DEL command",
			query:     "DEL user_****",
			mockCmd:   &parser.Command{CommandType: parser.DelCommandType, Arguments: []string{"user_****"}},
			mockErr:   nil,
			wantRes:   &Result{Command: parser.Command{CommandType: parser.DelCommandType, Arguments: []string{"user_****"}}},
			wantError: false,
		},
		{
			name:      "parser error",
			query:     "UPDATE something",
			mockCmd:   nil,
			mockErr:   errors.New("unknown command"),
			wantRes:   nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := parser.NewMockParser(t)
			mockParser.EXPECT().Parse(tt.query).Return(tt.mockCmd, tt.mockErr)

			interp := NewInterpreter(mockParser)
			res, err := interp.Interpret(tt.query)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if res == nil || !reflect.DeepEqual(res.Command, tt.wantRes.Command) {
					t.Errorf("got result %+v, want %+v", res, tt.wantRes)
				}
			}
		})
	}
}
