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

func TestParseResource(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantType BoxResourceType
		wantId   string
		wantErr  bool
	}{
		{"file", `{"type":"file","id":"f1","name":"a.txt"}`, FileResource, "f1", false},
		{"folder", `{"type":"folder","id":"d1","name":"dir"}`, FolderResource, "d1", false},
		{"file_version", `{"type":"file_version","id":"v1","sha1":"abc"}`, FileVersionResource, "v1", false},
		{"user", `{"type":"user","id":"u1","login":"u1@example.com"}`, UserResource, "u1", false},
		{"group", `{"type":"group","id":"g1","name":"grp"}`, GroupResource, "g1", false},
		{"group_membership", `{"type":"group_membership","id":"m1","role":"member"}`, MembershipResource, "m1", false},
		{"collaboration", `{"type":"collaboration","id":"c1","role":"editor"}`, CollaborationResource, "c1", false},
		{
			"top-level type wins over nested type that appears first",
			`{"parent":{"type":"folder","id":"d1"},"type":"file","id":"f2"}`,
			FileResource, "f2", false,
		},
		{"unsupported type", `{"type":"web_link","id":"w1"}`, 0, "", true},
		{"missing type", `{"id":"x1","name":"no type"}`, 0, "", true},
		{"invalid json", `{"type":"file",`, 0, "", true},
		{"empty", ``, 0, "", true},
		{"json null", `null`, 0, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseResource([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseResource() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != nil {
					t.Errorf("ParseResource() = %v, want nil on error", got)
				}
				return
			}
			if got.ResourceType() != tt.wantType {
				t.Errorf("ResourceType() = %v, want %v", got.ResourceType(), tt.wantType)
			}
			var id *string
			switch r := got.(type) {
			case *File:
				id = r.ID
			case *Folder:
				id = r.ID
			case *FileVersion:
				id = &r.ID
			case *User:
				id = r.ID
			case *Group:
				id = r.ID
			case *Membership:
				id = r.ID
			case *Collaboration:
				id = r.ID
			}
			if id == nil || *id != tt.wantId {
				t.Errorf("ID = %v, want %q", id, tt.wantId)
			}
		})
	}
}

func TestParseResource_DecodeError(t *testing.T) {
	// known type, but a field with the wrong JSON type: the partially decoded resource is returned with the error
	got, err := ParseResource([]byte(`{"type":"folder","id":"d1","name":123}`))
	if err == nil {
		t.Fatalf("ParseResource() expected a decode error")
	}
	folder, ok := got.(*Folder)
	if !ok {
		t.Fatalf("ParseResource() = %T, want *Folder together with the error", got)
	}
	if folder.ID == nil || *folder.ID != "d1" {
		t.Errorf("partially decoded ID = %v, want d1", folder.ID)
	}
}
