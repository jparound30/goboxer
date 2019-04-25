package goboxer

import (
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
		api *ApiConn
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
	type args struct {
		fields []string
		notify bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			if got := c.CreateReq(tt.args.fields, tt.args.notify); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.CreateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Create(t *testing.T) {
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
	type args struct {
		fields []string
		notify bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collaboration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			got, err := c.Create(tt.args.fields, tt.args.notify)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_UpdateReq(t *testing.T) {
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
	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			if got := c.UpdateReq(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.UpdateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Update(t *testing.T) {
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
	type args struct {
		collaborationId string
		role            Role
		status          *CollaborationStatus
		canViewPath     *bool
		fields          []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Collaboration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			got, err := c.Update(tt.args.collaborationId, tt.args.role, tt.args.status, tt.args.canViewPath, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_DeleteReq(t *testing.T) {
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
	type args struct {
		collaborationId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			if got := c.DeleteReq(tt.args.collaborationId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.DeleteReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_Delete(t *testing.T) {
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
	type args struct {
		collaborationId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			if err := c.Delete(tt.args.collaborationId); (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollaboration_PendingCollaborationsReq(t *testing.T) {
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
	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			if got := c.PendingCollaborationsReq(tt.args.offset, tt.args.limit, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collaboration.PendingCollaborationsReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollaboration_PendingCollaborations(t *testing.T) {
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
	type args struct {
		offset int
		limit  int
		fields []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantPendingList   []*Collaboration
		wantOutOffset     int
		wantOutLimit      int
		wantOutTotalCount int
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collaboration{
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
			gotPendingList, gotOutOffset, gotOutLimit, gotOutTotalCount, err := c.PendingCollaborations(tt.args.offset, tt.args.limit, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collaboration.PendingCollaborations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPendingList, tt.wantPendingList) {
				t.Errorf("Collaboration.PendingCollaborations() gotPendingList = %v, want %v", gotPendingList, tt.wantPendingList)
			}
			if gotOutOffset != tt.wantOutOffset {
				t.Errorf("Collaboration.PendingCollaborations() gotOutOffset = %v, want %v", gotOutOffset, tt.wantOutOffset)
			}
			if gotOutLimit != tt.wantOutLimit {
				t.Errorf("Collaboration.PendingCollaborations() gotOutLimit = %v, want %v", gotOutLimit, tt.wantOutLimit)
			}
			if gotOutTotalCount != tt.wantOutTotalCount {
				t.Errorf("Collaboration.PendingCollaborations() gotOutTotalCount = %v, want %v", gotOutTotalCount, tt.wantOutTotalCount)
			}
		})
	}
}
