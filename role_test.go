package goboxer

import (
	"reflect"
	"testing"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		name string
		r    *Role
		want string
	}{
		{"nil", nil, "<nil>"},
		{"normal", setRole(OWNER), "owner"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("Role.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       *Role
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte("null"), false},
		{"normal", setRole(CO_OWNER), []byte(`"co-owner"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Role.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
