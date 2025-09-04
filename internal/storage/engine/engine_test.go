package engine

import (
	"testing"
)

func TestEngine_SetGetDel(t *testing.T) {
	type op struct {
		action string
		key    string
		value  string
	}
	tests := []struct {
		name    string
		ops     []op
		wantGet map[string]struct {
			val string
			ok  bool
		}
	}{
		{
			name: "set and get single key",
			ops: []op{
				{action: "set", key: "foo", value: "bar"},
			},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"foo": {val: "bar", ok: true},
			},
		},
		{
			name: "set, get, del key",
			ops: []op{
				{action: "set", key: "foo", value: "bar"},
				{action: "del", key: "foo"},
			},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"foo": {val: "", ok: false},
			},
		},
		{
			name: "set multiple keys and get",
			ops: []op{
				{action: "set", key: "a", value: "1"},
				{action: "set", key: "b", value: "2"},
				{action: "set", key: "c", value: "3"},
			},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"a": {val: "1", ok: true},
				"b": {val: "2", ok: true},
				"c": {val: "3", ok: true},
			},
		},
		{
			name: "get non-existent key",
			ops:  []op{},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"missing": {val: "", ok: false},
			},
		},
		{
			name: "overwrite key",
			ops: []op{
				{action: "set", key: "foo", value: "bar"},
				{action: "set", key: "foo", value: "baz"},
			},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"foo": {val: "baz", ok: true},
			},
		},
		{
			name: "del non-existent key",
			ops: []op{
				{action: "del", key: "ghost"},
			},
			wantGet: map[string]struct {
				val string
				ok  bool
			}{
				"ghost": {val: "", ok: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEngine()
			for _, op := range tt.ops {
				switch op.action {
				case "set":
					e.Set(op.key, op.value)
				case "del":
					e.Del(op.key)
				}
			}
			for k, want := range tt.wantGet {
				gotVal, gotOk := e.Get(k)
				if gotVal != want.val || gotOk != want.ok {
					t.Errorf("Get(%q) = (%q, %v), want (%q, %v)", k, gotVal, gotOk, want.val, want.ok)
				}
			}
		})
	}
}
