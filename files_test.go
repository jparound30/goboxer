package goboxer

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io/ioutil"
	"testing"
	"time"
)

func TestFile_Unmarshal(t *testing.T) {
	typ := TYPE_FILE
	id := "5000948880"
	fileVersion := FileVersion{Type: "file_version", ID: "26261748416", Sha1: "134b65991ed521fcfe4724b7d814ab8ded5185dc"}
	sequenceId := "3"
	etag := "3"
	sha1 := "134b65991ed521fcfe4724b7d814ab8ded5185dc"
	name := "tigers.jpeg"
	description := "a picture of tigers"
	size := float64(629644)

	item1Type := TYPE_FOLDER
	item1Id := "0"
	item1Name := "All Files"
	path1 := &ItemMini{Type: &item1Type, ID: &item1Id, SequenceId: nil, ETag: nil, Name: &item1Name}

	item2Type := TYPE_FOLDER
	item2Id := "11446498"
	item2Name := "Pictures"
	item2Seq := "1"
	item2Etag := "1"
	path2 := &ItemMini{Type: &item2Type, ID: &item2Id, SequenceId: &item2Seq, ETag: &item2Etag, Name: &item2Name}
	pathCollection := PathCollection{
		TotalCount: 2,
		Entries:    []*ItemMini{path1, path2},
	}

	createdAt, _ := time.Parse(time.RFC3339, "2012-12-12T10:55:30-08:00")
	modifiedAt, _ := time.Parse(time.RFC3339, "2012-12-12T11:04:26-08:00")

	contentCreatedAt, _ := time.Parse(time.RFC3339, "2013-02-04T16:57:52-08:00")
	contentModifiedAt, _ := time.Parse(time.RFC3339, "2013-02-04T16:57:52-08:00")

	ctype := TYPE_USER
	cid := "17738362"
	cname := "sean rose"
	clogin := "sean@box.com"
	createdBy := UserGroupMini{Type: &ctype, ID: &cid, Name: &cname, Login: &clogin}

	mtype := TYPE_USER
	mid := "17738362"
	mname := "sean rose"
	mlogin := "sean@box.com"
	modifiedBy := UserGroupMini{Type: &mtype, ID: &mid, Name: &mname, Login: &mlogin}

	otype := TYPE_USER
	oid := "17738362"
	oname := "sean rose"
	ologin := "sean@box.com"
	ownedBy := UserGroupMini{Type: &otype, ID: &oid, Name: &oname, Login: &ologin}

	sharedLink := &SharedLink{}
	slUrl := "https://www.box.com/s/rh935iit6ewrmw0unyul"
	sharedLink.Url = &slUrl
	slDownloadUrl := "https://www.box.com/shared/static/rh935iit6ewrmw0unyul.jpeg"
	sharedLink.DownloadUrl = &slDownloadUrl
	sharedLink.VanityUrl = nil
	slIsPasswordEnabled := false
	sharedLink.IsPasswordEnabled = &slIsPasswordEnabled
	sharedLink.UnsharedAt = nil
	slDownloadCount := 0
	sharedLink.DownloadCount = &slDownloadCount
	slPreviewCount := 0
	sharedLink.PreviewCount = &slPreviewCount
	slAccess := "open"
	sharedLink.Access = &slAccess
	tr := true
	sharedLink.Permissions = &Permissions{
		CanDownload: &tr,
		CanPreview:  &tr,
	}

	ptype := TYPE_FOLDER
	pid := "11446498"
	pseqid := "1"
	petag := "1"
	pname := "Pictures"
	parent := ItemMini{
		Type:       &ptype,
		ID:         &pid,
		SequenceId: &pseqid,
		ETag:       &petag,
		Name:       &pname,
	}
	itemStatus := "active"
	tags := []string{"cropped", "color corrected"}
	ltype := "lock"
	lid := "112429"
	lctype := TYPE_USER
	lcid := "18212074"
	lcname := "Obi Wan"
	lclogin := "obiwan@box.com"
	lcreatedBy := UserGroupMini{
		&lctype,
		&lcid,
		&lcname,
		&lclogin,
	}
	lcreatedAt, _ := time.Parse(time.RFC3339, "2013-12-04T10:28:36-08:00")
	lexpiresAt, _ := time.Parse(time.RFC3339, "2012-12-12T10:55:30-08:00")
	lisDownloadPrevented := false
	lock := Lock{
		Type:                &ltype,
		ID:                  &lid,
		CreatedBy:           &lcreatedBy,
		CreatedAt:           &lcreatedAt,
		ExpiresAt:           &lexpiresAt,
		IsDownloadPrevented: &lisDownloadPrevented,
	}
	tests := []struct {
		name     string
		jsonFile string
		want     File
	}{
		// TODO: Add test cases.
		{
			name:     "normal",
			jsonFile: "testdata/files/file_json.json",
			want: File{
				ItemMini{
					Type:       &typ,
					ID:         &id,
					SequenceId: &sequenceId,
					ETag:       &etag,
					Name:       &name,
				},
				nil,
				&fileVersion,
				&sha1,
				&description,
				size,
				&pathCollection,
				&createdAt,
				&modifiedAt,
				nil,
				nil,
				&contentCreatedAt,
				&contentModifiedAt,
				nil,
				&createdBy,
				&modifiedBy,
				&ownedBy,
				sharedLink,
				&parent,
				&itemStatus,
				nil,
				nil,
				nil,
				tags,
				&lock,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := ioutil.ReadFile(tt.jsonFile)
			file := File{}
			err := json.Unmarshal(b, &file)
			if err != nil {
				t.Errorf("File Unmarshal err %v", err)
			}
			opt := cmpopts.IgnoreUnexported(file)
			if diff := cmp.Diff(&file, &tt.want, opt); diff != "" {
				t.Errorf("File Marshal/Unmarshal differs: (-got +want)\n%s", diff)
			}
		})
	}
}
