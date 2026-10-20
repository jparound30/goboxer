package goboxer

import "testing"

func Test_intToString(t *testing.T) {
	tests := []struct {
		name string
		in   *int
		want string
	}{
		{"nil", nil, "<nil>"},
		{"zero", setIntPtr(0), "0"},
		{"positive", setIntPtr(65), "65"},
		{"large", setIntPtr(123456), "123456"},
		{"negative", setIntPtr(-7), "-7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intToString(tt.in); got != tt.want {
				t.Errorf("intToString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSharedLink_String(t *testing.T) {
	var nilLink *SharedLink
	if got := nilLink.String(); got != "<nil>" {
		t.Errorf("nil SharedLink.String() = %q, want %q", got, "<nil>")
	}

	sl := &SharedLink{
		Url:               setStringPtr("https://app.box.com/s/abc"),
		IsPasswordEnabled: setBool(false),
		DownloadCount:     setIntPtr(65),
		PreviewCount:      setIntPtr(1200),
		Access:            setStringPtr("open"),
	}
	want := "{ Url:https://app.box.com/s/abc, DownloadUrl:<nil>, VanityUrl:<nil>, IsPasswordEnabled:false, UnsharedAt:<nil>," +
		" DownloadCount:65, PreviewCount:1200, Access:open, Permissions:<nil>, Password:<nil>}"
	if got := sl.String(); got != want {
		t.Errorf("SharedLink.String() =\n%s\nwant\n%s", got, want)
	}
}
