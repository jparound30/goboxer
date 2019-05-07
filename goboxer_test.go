package goboxer

import (
	"reflect"
	"testing"
)

func TestItemType_String(t *testing.T) {
	tests := []struct {
		name string
		i    *ItemType
		want string
	}{
		{"nil", nil, "<nil>"},
		{"normal", setItemTypePtr(TYPE_FILE), "file"},
		{"normal", setItemTypePtr(TYPE_FOLDER), "folder"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("ItemType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		i       *ItemType
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"normal", setItemTypePtr(TYPE_FILE), []byte(`"file"`), false},
		{"normal", setItemTypePtr(TYPE_FOLDER), []byte(`"folder"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ItemType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
