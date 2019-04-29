package goboxer

import "testing"

func TestFileVersion_ResourceType(t *testing.T) {
	tests := []struct {
		name   string
		target *FileVersion
		want   BoxResourceType
	}{
		{"nil", nil, FileVersionResource},
		{"normal", &FileVersion{}, FileVersionResource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := tt.target
			if got := fv.ResourceType(); got != tt.want {
				t.Errorf("FileVersion.ResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}
