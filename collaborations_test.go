package goboxer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func buildCollaborationsOfGetInfoNormalJson() *Collaboration {
	var n Collaboration
	n.Type = setStringPtr("collaboration")
	n.ID = setStringPtr("791293")

	n.CreatedBy = &UserGroupMini{
		Type:  setUserType(TYPE_USER),
		ID:    setStringPtr("17738362"),
		Name:  setStringPtr("sean rose"),
		Login: setStringPtr("sean@box.com"),
	}
	n.CreatedAt = setTime("2012-12-12T10:54:37-08:00")
	n.ModifiedAt = setTime("2012-12-12T11:30:43-08:00")
	n.ExpiresAt = nil
	n.Status = setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED)
	n.AccessibleBy = &UserGroupMini{
		Type:  setUserType(TYPE_USER),
		ID:    setStringPtr("18203124"),
		Name:  setStringPtr("sean"),
		Login: setStringPtr("sean+test@box.com"),
	}
	n.Role = setRole(EDITOR)
	n.AcknowledgedAt = setTime("2012-12-12T11:30:43-08:00")
	n.Item = &ItemMini{
		Type:       setItemTypePtr(TYPE_FOLDER),
		ID:         setStringPtr("11446500"),
		SequenceId: setStringPtr("0"),
		ETag:       setStringPtr("0"),
		Name:       setStringPtr("Shared Pictures"),
	}
	return &n
}

func TestCollaborationStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		cs      *CollaborationStatus
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"accepted", setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED), []byte(`"accepted"`), false},
		{"pending", setCollaborationStatus(COLLABORATION_STATUS_PENDING), []byte(`"pending"`), false},
		{"rejected", setCollaborationStatus(COLLABORATION_STATUS_REJECTED), []byte(`"rejected"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cs.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("CollaborationStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollaborationStatus.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaborationStatus_String(t *testing.T) {
	tests := []struct {
		name string
		cs   *CollaborationStatus
		want string
	}{
		{"nil", nil, "<nil>"},
		{"accepted", setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED), "accepted"},
		{"pending", setCollaborationStatus(COLLABORATION_STATUS_PENDING), "pending"},
		{"rejected", setCollaborationStatus(COLLABORATION_STATUS_REJECTED), "rejected"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.String(); got != tt.want {
				t.Errorf("CollaborationStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_String(t *testing.T) {
	type fields struct {
		apiInfo        *apiInfo
		Type           *string
		ID             *string
		CreatedBy      *UserGroupMini
		CreatedAt      *time.Time
		ModifiedAt     *time.Time
		ExpiresAt      *time.Time
		Status         *CollaborationStatus
		AccessibleBy   *UserGroupMini
		InviteEmail    *string
		Role           *Role
		AcknowledgedAt *time.Time
		Item           *ItemMini
		CanViewPath    *bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"nil",
			fields{
				Type: nil,
			},
			"<nil>",
		},
		{"normal",
			fields{
				Type: setStringPtr("collaboration"),
				ID:   setStringPtr("791293"),
				CreatedBy: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("17738362"),
					Name:  setStringPtr("sean rose"),
					Login: setStringPtr("sean@box.com"),
				},
				CreatedAt:  setTime("2012-12-12T10:54:37-08:00"),
				ModifiedAt: setTime("2012-12-12T11:30:43-08:00"),
				ExpiresAt:  nil,
				Status:     setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED),
				AccessibleBy: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("18203124"),
					Name:  setStringPtr("sean"),
					Login: setStringPtr("sean+test@box.com"),
				},
				Role:           setRole(EDITOR),
				AcknowledgedAt: setTime("2012-12-12T11:30:43-08:00"),
				Item: &ItemMini{
					Type:       setItemTypePtr(TYPE_FOLDER),
					ID:         setStringPtr("11446500"),
					SequenceId: setStringPtr("0"),
					ETag:       setStringPtr("0"),
					Name:       setStringPtr("Shared Pictures"),
				},
			},
			"{ Type:collaboration, ID:791293, CreatedBy:{ Type:user, ID:17738362, Name:sean rose, Login:sean@box.com }, " +
				"CreatedAt:2012-12-12T10:54:37-08:00, Modified:2012-12-12T11:30:43-08:00, ExpiresAt:<nil>, Status:accepted," +
				" AccessibleBy:{ Type:user, ID:18203124, Name:sean, Login:sean+test@box.com }, InviteEmail:<nil>, Role:editor, " +
				"AcknowledgedAt:2012-12-12T11:30:43-08:00, Item:{ Type:folder, ID:11446500, Name:Shared Pictures, SequenceId:0, ETag:0 }, " +
				"CanViewPath:<nil> }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *Collaboration
			if tt.fields.Type == nil {
				c = nil
			} else {
				c = &Collaboration{
					apiInfo:        tt.fields.apiInfo,
					Type:           tt.fields.Type,
					ID:             tt.fields.ID,
					CreatedBy:      tt.fields.CreatedBy,
					CreatedAt:      tt.fields.CreatedAt,
					ModifiedAt:     tt.fields.ModifiedAt,
					ExpiresAt:      tt.fields.ExpiresAt,
					Status:         tt.fields.Status,
					AccessibleBy:   tt.fields.AccessibleBy,
					InviteEmail:    tt.fields.InviteEmail,
					Role:           tt.fields.Role,
					AcknowledgedAt: tt.fields.AcknowledgedAt,
					Item:           tt.fields.Item,
					CanViewPath:    tt.fields.CanViewPath,
				}
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Collaboration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_ResourceType(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type fields struct {
		target *Collaboration
	}
	tests := []struct {
		name   string
		fields fields
		want   BoxResourceType
	}{
		{"nil", fields{nil}, CollaborationResource},
		{"normal", fields{&Collaboration{apiInfo: &apiInfo{api: apiConn}}}, CollaborationResource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.target.ResourceType(); got != tt.want {
				t.Errorf("Collaboration.ResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCollaboration(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		api *APIConn
	}
	tests := []struct {
		name string
		args args
		want *Collaboration
	}{
		{"normal", args{api: apiConn}, &Collaboration{apiInfo: &apiInfo{api: apiConn}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCollaboration(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCollaboration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_GetInfoReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		collaborationId string
		fields          []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{"123", nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/collaborations/123",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{"123", []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/collaborations/123?fields=type,id",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollaboration(apiConn)
			got := c.GetInfoReq(tt.args.collaborationId, tt.args.fields)
			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestCollaboration_GetInfo(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/collaborations/") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/collaborations/")
			}
			// Method check
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			collabId := strings.TrimPrefix(r.URL.Path, "/2.0/collaborations/")

			switch collabId {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				resp, _ := ioutil.ReadFile("testdata/collaborations/collaboration_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildCollaborationsOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		collaborationId string
		fields          []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Collaboration
		wantErr bool
		errType interface{}
	}{
		{"normal/fields unspecified", args{collaborationId: "10001", fields: nil}, normal, false, nil},
		{"normal/allFields", args{collaborationId: "10002", fields: CollaborationAllFields}, normal, false, nil},
		{"http error/404", args{collaborationId: "404", fields: CollaborationAllFields}, nil, true, &ApiStatusError{}},
		{"returned invalid json/999", args{collaborationId: "999", fields: nil}, nil, true, &ApiOtherError{}},
		{"senderror", args{collaborationId: "999", fields: nil}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			c := NewCollaboration(apiConn)
			got, err := c.GetInfo(tt.args.collaborationId, tt.args.fields)
			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					if status, err := strconv.Atoi(tt.args.collaborationId); err != nil || status != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			// If normal response
			opt := cmpopts.IgnoreUnexported(*got, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want, opt); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			if got.apiInfo == nil {
				t.Errorf("not exists `apiInfo` field\n")
				return
			}
		})
	}
}

func TestCollaboration_SetItem(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1before := buildCollaborationsOfGetInfoNormalJson()
	c1before.apiInfo = &apiInfo{api: apiConn}

	c1after := buildCollaborationsOfGetInfoNormalJson()
	c1after.apiInfo = &apiInfo{api: apiConn}
	c1after.Item = &ItemMini{
		Type: setItemTypePtr(TYPE_FOLDER),
		ID:   setStringPtr("ID1"),
	}

	type args struct {
		typ ItemType
		id  string
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Collaboration
	}{
		{"normal", c1before, args{typ: TYPE_FOLDER, id: "ID1"}, c1after},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			got := c.SetItem(tt.args.typ, tt.args.id)
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_SetAccessibleById(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1before := buildCollaborationsOfGetInfoNormalJson()
	c1before.apiInfo = &apiInfo{api: apiConn}

	c1after := buildCollaborationsOfGetInfoNormalJson()
	c1after.apiInfo = &apiInfo{api: apiConn}
	c1after.AccessibleBy = &UserGroupMini{
		Type: setUserType(TYPE_GROUP),
		ID:   setStringPtr("ID1"),
	}

	type args struct {
		typ UserGroupType
		id  string
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Collaboration
	}{
		{"normal", c1before, args{typ: TYPE_GROUP, id: "ID1"}, c1after},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			got := c.SetAccessibleById(tt.args.typ, tt.args.id)
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_SetAccessibleByEmailForNewUser(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1before := buildCollaborationsOfGetInfoNormalJson()
	c1before.apiInfo = &apiInfo{api: apiConn}

	c1after := buildCollaborationsOfGetInfoNormalJson()
	c1after.apiInfo = &apiInfo{api: apiConn}
	c1after.AccessibleBy = &UserGroupMini{
		Type:  setUserType(TYPE_USER),
		Login: setStringPtr("LOGIN1@example.com"),
	}

	type args struct {
		login string
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Collaboration
	}{
		{"normal", c1before, args{login: "LOGIN1@example.com"}, c1after},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			got := c.SetAccessibleByEmailForNewUser(tt.args.login)
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_SetRole(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1before := buildCollaborationsOfGetInfoNormalJson()
	c1before.apiInfo = &apiInfo{api: apiConn}

	c1after := buildCollaborationsOfGetInfoNormalJson()
	c1after.apiInfo = &apiInfo{api: apiConn}
	c1after.Role = setRole(CO_OWNER)

	type args struct {
		role Role
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Collaboration
	}{
		{"normal", c1before, args{CO_OWNER}, c1after},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			got := c.SetRole(tt.args.role)
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_SetCanViewPath(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1before := buildCollaborationsOfGetInfoNormalJson()
	c1before.apiInfo = &apiInfo{api: apiConn}

	c1after := buildCollaborationsOfGetInfoNormalJson()
	c1after.apiInfo = &apiInfo{api: apiConn}
	c1after.CanViewPath = setBool(true)

	type args struct {
		canViewPath bool
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Collaboration
	}{
		{"normal", c1before, args{true}, c1after},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target
			got := c.SetCanViewPath(tt.args.canViewPath)
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_CreateReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1 := buildCollaborationsOfGetInfoNormalJson()
	c1.apiInfo = &apiInfo{api: apiConn}
	w1 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "folder",
		"id": "FOLDER_ID1"
	},
	"accessible_by": {
		"type": "user",
		"id": "USER_ID1"
	},
	"role": "editor"
}
`),
	}

	c2 := buildCollaborationsOfGetInfoNormalJson()
	c2.apiInfo = &apiInfo{api: apiConn}
	w2 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "folder",
		"id": "FOLDER_ID1"
	},
	"accessible_by": {
		"type": "group",
		"id": "GROUP_ID1"
	},
	"role": "co-owner"
}
`),
	}

	c3 := buildCollaborationsOfGetInfoNormalJson()
	c3.apiInfo = &apiInfo{api: apiConn}
	w3 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "folder",
		"id": "FOLDER_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "viewer"
}
`),
	}

	c4 := buildCollaborationsOfGetInfoNormalJson()
	c4.apiInfo = &apiInfo{api: apiConn}
	w4 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"id": "USER_ID1"
	},
	"role": "editor"
}
`),
	}

	c5 := buildCollaborationsOfGetInfoNormalJson()
	c5.apiInfo = &apiInfo{api: apiConn}
	w5 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "group",
		"id": "GROUP_ID1"
	},
	"role": "previewer"
}
`),
	}

	c6 := buildCollaborationsOfGetInfoNormalJson()
	c6.apiInfo = &apiInfo{api: apiConn}
	w6 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "previewer uploader"
}
`),
	}

	c7 := buildCollaborationsOfGetInfoNormalJson()
	c7.apiInfo = &apiInfo{api: apiConn}
	w7 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "viewer uploader",
	"can_view_path": true
}
`),
	}

	c8 := c7
	w8 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=false&fields=type,id",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "viewer uploader",
	"can_view_path": true
}
`),
	}

	c9 := c7
	w9 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=true",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "viewer uploader",
	"can_view_path": true
}
`),
	}

	c10 := c7
	w10 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations?notify=true&fields=type,id",
		Method:             POST,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"item": {
		"type": "file",
		"id": "FILE_ID1"
	},
	"accessible_by": {
		"type": "user",
		"login": "BBB@example.com"
	},
	"role": "viewer uploader",
	"can_view_path": true
}
`),
	}

	type args struct {
		targetItem  ItemMini
		grantedTo   UserGroupMini
		role        Role
		canViewPath *bool
		fields      []string
		notify      bool
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Request
	}{
		{"normal/folder/granted to user",
			c1,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("FOLDER_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("USER_ID1")},
				EDITOR,
				nil,
				nil,
				false,
			},
			w1,
		},
		{"normal/folder/granted to group",
			c2,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("FOLDER_ID1")},
				UserGroupMini{Type: setUserType(TYPE_GROUP), ID: setStringPtr("GROUP_ID1")},
				CO_OWNER,
				nil,
				nil,
				false,
			},
			w2,
		},
		{"normal/folder/granted to new user",
			c3,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FOLDER), ID: setStringPtr("FOLDER_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER,
				nil,
				nil,
				false,
			},
			w3,
		},
		{"normal/file/granted to user",
			c4,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), ID: setStringPtr("USER_ID1")},
				EDITOR,
				nil,
				nil,
				false,
			},
			w4,
		},
		{"normal/file/granted to group",
			c5,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_GROUP), ID: setStringPtr("GROUP_ID1")},
				PREVIEWER,
				nil,
				nil,
				false,
			},
			w5,
		},
		{"normal/file/granted to new user",
			c6,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				PREVIEWER_UPLOADER,
				nil,
				nil,
				false,
			},
			w6,
		},
		{"normal/can view path",
			c7,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				nil,
				false,
			},
			w7,
		},
		{"normal/fields",
			c8,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				false,
			},
			w8,
		},
		{"normal/notify",
			c9,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				nil,
				true,
			},
			w9,
		},
		{"normal/fields and notify",
			c10,
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				true,
			},
			w10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			c := tt.target
			got := c.CreateReq(tt.args.targetItem, tt.args.grantedTo, tt.args.role, tt.args.canViewPath, tt.args.fields, tt.args.notify)

			opts := diffCompOptions(Collaboration{}, APIConn{})
			opt := cmpopts.IgnoreUnexported(Request{})
			opts = append(opts, opt)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
			gotBodyDec := json.NewDecoder(got.body)
			var gotBody map[string]interface{}
			err := gotBodyDec.Decode(&gotBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			expBodyDec := json.NewDecoder(tt.want.body)
			var expBody map[string]interface{}
			err = expBodyDec.Decode(&expBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			if diff := cmp.Diff(gotBody, expBody); diff != "" {
				t.Errorf("body differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_Create(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/collaborations") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/collaborations")
			}
			// Method check
			if r.Method != http.MethodPost {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			body, _ := ioutil.ReadAll(r.Body)
			var js map[string]interface{}
			_ = json.Unmarshal(body, &js)
			tmp1 := js["item"].(map[string]interface{})
			collabId := tmp1["id"].(string)

			switch collabId {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(201)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(201)
				resp, _ := ioutil.ReadFile("testdata/collaborations/collaboration_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildCollaborationsOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		targetItem  ItemMini
		grantedTo   UserGroupMini
		role        Role
		canViewPath *bool
		fields      []string
		notify      bool
	}
	tests := []struct {
		name    string
		args    args
		want    *Collaboration
		wantErr bool
		errType interface{}
	}{
		{"normal",
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("FILE_ID1")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				true,
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("404")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				true,
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("999")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				true,
			}, nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			args{
				ItemMini{Type: setItemTypePtr(TYPE_FILE), ID: setStringPtr("999")},
				UserGroupMini{Type: setUserType(TYPE_USER), Login: setStringPtr("BBB@example.com")},
				VIEWER_UPLOADER,
				setBool(true),
				[]string{"type", "id"},
				true,
			},
			nil,
			true,
			&ApiOtherError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			c := NewCollaboration(apiConn)
			got, err := c.Create(tt.args.targetItem, tt.args.grantedTo, tt.args.role, tt.args.canViewPath, tt.args.fields, tt.args.notify)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(*ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			// If normal response
			opt := cmpopts.IgnoreUnexported(*got, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want, opt); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			if got.apiInfo == nil {
				t.Errorf("not exists `apiInfo` field\n")
				return
			}
		})
	}
}

func TestCollaboration_UpdateReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	c1 := buildCollaborationsOfGetInfoNormalJson()
	c1.apiInfo = &apiInfo{api: apiConn}

	w1 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001",
		Method:             PUT,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"role": "editor"
}
`),
	}

	w2 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001",
		Method:             PUT,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"role": "viewer",
	"status": "accepted"
}
`),
	}

	w3 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001",
		Method:             PUT,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"role": "previewer",
	"status": "rejected",
	"can_view_path": true
}
`),
	}

	w4 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001",
		Method:             PUT,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"role": "uploader",
	"status": "pending",
	"can_view_path": false
}
`),
	}

	w5 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001?fields=type,id",
		Method:             PUT,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body: strings.NewReader(`
{
	"role": "previewer uploader",
	"status": "pending",
	"can_view_path": false
}
`),
	}

	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name   string
		target *Collaboration
		args   args
		want   *Request
	}{
		{"role",
			c1,
			args{
				collaborationId: "10001",
				role:            EDITOR,
				status:          nil,
				canViewPath:     nil,
				fields:          nil,
			},
			w1,
		},
		{"role, status",
			c1,
			args{
				collaborationId: "10001",
				role:            VIEWER,
				status:          setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED),
				canViewPath:     nil,
				fields:          nil,
			},
			w2,
		},
		{"role, status,can view path",
			c1,
			args{
				collaborationId: "10001",
				role:            PREVIEWER,
				status:          setCollaborationStatus(COLLABORATION_STATUS_REJECTED),
				canViewPath:     setBool(true),
				fields:          nil,
			},
			w3,
		},
		{"role, status,can view path",
			c1,
			args{
				collaborationId: "10001",
				role:            UPLOADER,
				status:          setCollaborationStatus(COLLABORATION_STATUS_PENDING),
				canViewPath:     setBool(false),
				fields:          nil,
			},
			w4,
		},
		{"role, status,can view path, fields",
			c1,
			args{
				collaborationId: "10001",
				role:            PREVIEWER_UPLOADER,
				status:          setCollaborationStatus(COLLABORATION_STATUS_PENDING),
				canViewPath:     setBool(false),
				fields:          []string{"type", "id"},
			},
			w5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.target

			got := c.UpdateReq(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields)

			opts := diffCompOptions(Collaboration{}, APIConn{})
			opt := cmpopts.IgnoreUnexported(Request{})
			opts = append(opts, opt)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
			gotBodyDec := json.NewDecoder(got.body)
			var gotBody map[string]interface{}
			err := gotBodyDec.Decode(&gotBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			expBodyDec := json.NewDecoder(tt.want.body)
			var expBody map[string]interface{}
			err = expBodyDec.Decode(&expBody)
			if err != nil {
				t.Fatalf("body json doesnt unmarshal")
			}

			if diff := cmp.Diff(gotBody, expBody); diff != "" {
				t.Errorf("body differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_Update(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/collaborations") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/collaborations")
			}
			// Method check
			if r.Method != http.MethodPut {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			collabId := strings.Split(r.URL.Path, "/")[3]

			switch collabId {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			case "10002":
				w.WriteHeader(204)
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/collaborations/collaboration_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildCollaborationsOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Collaboration
		wantErr bool
		errType interface{}
	}{
		{"normal",
			args{
				"10001",
				VIEWER_UPLOADER,
				nil,
				nil,
				[]string{"type", "id"},
			},
			normal,
			false,
			nil,
		},
		{"normal",
			args{
				"10002",
				VIEWER_UPLOADER,
				nil,
				nil,
				[]string{"type", "id"},
			},
			nil,
			false,
			nil,
		},
		{"http error/404",
			args{
				"404",
				VIEWER_UPLOADER,
				nil,
				nil,
				[]string{"type", "id"},
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			args{
				"999",
				VIEWER_UPLOADER,
				nil,
				nil,
				[]string{"type", "id"},
			},
			normal,
			true,
			&ApiOtherError{}},
		{"senderror",
			args{
				"999",
				VIEWER_UPLOADER,
				nil,
				nil,
				[]string{"type", "id"},
			},
			nil,
			true,
			&ApiOtherError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			c := NewCollaboration(apiConn)

			got, err := c.Update(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(*ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}

			if got == nil {
				return
			}
			// If normal response
			opt := cmpopts.IgnoreUnexported(*got, Collaboration{})
			if diff := cmp.Diff(&got, &tt.want, opt); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			if got.apiInfo == nil {
				t.Errorf("not exists `apiInfo` field\n")
				return
			}
		})
	}
}

func TestCollaboration_DeleteReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	w1 := &Request{
		apiConn:            apiConn,
		Url:                url + "/2.0/collaborations/10001",
		Method:             DELETE,
		shouldAuthenticate: true,
		numRedirects:       defaultNumRedirects,
		headers:            http.Header{},
		body:               nil,
	}

	type args struct {
		collaborationId string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal", args{"10001"}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollaboration(apiConn)

			got := c.DeleteReq(tt.args.collaborationId)
			opts := diffCompOptions(Collaboration{}, APIConn{})
			opt := cmpopts.IgnoreUnexported(Request{})
			opts = append(opts, opt)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestCollaboration_Delete(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/collaborations") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/collaborations")
			}
			// Method check
			if r.Method != http.MethodDelete {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			collabId := strings.Split(r.URL.Path, "/")[3]

			switch collabId {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(204)
				_, _ = w.Write([]byte("invalid json"))
			case "10002":
				w.WriteHeader(204)
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(204)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildCollaborationsOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		collaborationId string
	}
	tests := []struct {
		name    string
		args    args
		want    *Collaboration
		wantErr bool
		errType interface{}
	}{
		{"normal",
			args{
				"10001",
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			args{
				"404",
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"senderror",
			args{
				"999",
			},
			nil,
			true,
			&ApiOtherError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			c := NewCollaboration(apiConn)
			err := c.Delete(tt.args.collaborationId)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(*ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}
		})
	}
}

func TestCollaboration_PendingCollaborationsReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{
			name: "normal/fields=nil",
			args: args{0, 1000, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/collaborations?status=pending&offset=0&limit=1000",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{2000, 1000, []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/collaborations?status=pending&offset=2000&limit=1000&fields=type,id",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields=nil",
			args: args{0, 1000, nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/collaborations?status=pending&offset=0&limit=1000",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollaboration(apiConn)

			got := c.PendingCollaborationsReq(tt.args.offset, tt.args.limit, tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got, Request{})
			opts = append(opts, cmpopts.IgnoreInterfaces(struct{ io.Reader }{}))
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestCollaboration_PendingCollaborations(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/collaborations") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/collaborations")
			}
			// Method check
			if r.Method != http.MethodGet {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			offset := r.URL.Query().Get("offset")
			limit := r.URL.Query().Get("limit")

			switch offset {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "999":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				base := `
{
  "entries":[
    {
      "type":"collaboration",
      "id":"27513888",
      "created_by":{
        "type":"user",
        "id":"11993747",
        "name":"sean",
        "login":"sean@box.com"
      },
      "created_at":"2012-10-17T23:14:42-07:00",
      "modified_at":"2012-10-17T23:14:42-07:00",
      "expires_at":null,
      "status":"pending",
      "accessible_by":{
        "type":"user",
        "id":"181216415",
        "name":"sean rose",
        "login":"sean+awesome@box.com"
      },
      "role":"editor",
      "acknowledged_at":null,
      "item":null
    }
  ],
	"total_count": 1000,
	"offset": $$$offset$$$,
	"limit": $$$limit$$$
}`
				base = strings.ReplaceAll(base, "$$$offset$$$", offset)
				base = strings.ReplaceAll(base, "$$$limit$$$", limit)
				_, _ = w.Write([]byte(base))
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildCollaborationsOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	entries1 := &Collaboration{
		apiInfo: &apiInfo{api: apiConn},
		Type:    setStringPtr("collaboration"),
		ID:      setStringPtr("27513888"),
		CreatedBy: &UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("11993747"),
			Name:  setStringPtr("sean"),
			Login: setStringPtr("sean@box.com"),
		},
		CreatedAt:  setTime("2012-10-17T23:14:42-07:00"),
		ModifiedAt: setTime("2012-10-17T23:14:42-07:00"),
		ExpiresAt:  nil,
		Status:     setCollaborationStatus(COLLABORATION_STATUS_PENDING),
		AccessibleBy: &UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("181216415"),
			Name:  setStringPtr("sean rose"),
			Login: setStringPtr("sean+awesome@box.com"),
		},
		Role:           setRole(EDITOR),
		AcknowledgedAt: nil,
		Item:           nil,
	}
	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name              string
		args              args
		wantPendingList   []*Collaboration
		wantOutOffset     int
		wantOutLimit      int
		wantOutTotalCount int
		wantErr           bool
		errType           interface{}
	}{
		{"normal",
			args{offset: 0, limit: 1000, fields: nil},
			[]*Collaboration{entries1},
			0,
			1000,
			1000,
			false,
			nil,
		},
		{"http error/404",
			args{offset: 404, limit: 1000, fields: nil},
			nil,
			404,
			1000,
			1000,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			args{
				999,
				1000,
				nil,
			},
			nil,
			999,
			1000,
			1000,
			true,
			&ApiOtherError{},
		},
		{"senderror",
			args{
				999,
				1000,
				nil,
			},
			nil,
			999,
			1000,
			1000,
			true,
			&ApiOtherError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			c := NewCollaboration(apiConn)
			gotPendingList, gotOutOffset, gotOutLimit, gotOutTotalCount, err := c.PendingCollaborations(tt.args.offset, tt.args.limit, tt.args.fields)

			// Error checks
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != nil {
				if reflect.TypeOf(err).String() != reflect.TypeOf(tt.errType).String() {
					t.Errorf("got err = %v, wanted errorType %v", err, tt.errType)
					return
				}
				if reflect.TypeOf(tt.errType) == reflect.TypeOf(&ApiStatusError{}) {
					apiStatusError := err.(*ApiStatusError)
					expectedStatus := tt.errType.(*ApiStatusError).Status
					if expectedStatus != apiStatusError.Status {
						t.Errorf("status code may be not corrected [%d]", apiStatusError.Status)
						return
					}
					return
				} else {
					return
				}
			} else if err != nil {
				return
			}
			opts := diffCompOptions(Collaboration{}, Request{}, apiInfo{})
			opts = append(opts, cmpopts.IgnoreInterfaces(struct{ io.Reader }{}))
			if diff := cmp.Diff(gotPendingList, tt.wantPendingList, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
			if gotOutOffset != tt.wantOutOffset {
				t.Errorf("invalid output pendingList")
			}
			if gotOutLimit != tt.wantOutLimit {
				t.Errorf("invalid output pendingList")
			}
			if gotOutTotalCount != tt.wantOutTotalCount {
				t.Errorf("invalid output pendingList")
			}
		})
	}
}
