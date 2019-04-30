package goboxer

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func setInvitabilityLevel(level InvitabilityLevel) *InvitabilityLevel {
	return &level
}
func setMemberViewabilityLevel(level MemberViewabilityLevel) *MemberViewabilityLevel {
	return &level
}

func buildGroupOfGetInfoNormalJson() *Group {
	var normal Group
	normal.Type = setUserType(TYPE_GROUP)
	normal.ID = setStringPtr("255224")
	normal.Name = setStringPtr("Everyone")
	normal.CreatedAt = setTime("2014-09-15T13:15:35-07:00")
	normal.ModifiedAt = setTime("2014-09-15T13:15:35-07:00")
	return &normal
}

func TestInvitabilityLevel_String(t *testing.T) {
	tests := []struct {
		name string
		us   *InvitabilityLevel
		want string
	}{
		// TODO: Add test cases.
		{"nil", nil, "<nil>"},
		{"normal", setInvitabilityLevel(InvitabilityAdminsOnly), "admins_only"},
		{"normal", setInvitabilityLevel(InvitabilityAdminsMembers), "admins_and_members"},
		{"normal", setInvitabilityLevel(InvitabilityAllManagedUsers), "all_managed_users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.String(); got != tt.want {
				t.Errorf("InvitabilityLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvitabilityLevel_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		us      *InvitabilityLevel
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte("null"), false},
		{"normal", setInvitabilityLevel(InvitabilityAdminsOnly), []byte(`"admins_only"`), false},
		{"normal", setInvitabilityLevel(InvitabilityAdminsMembers), []byte(`"admins_and_members"`), false},
		{"normal", setInvitabilityLevel(InvitabilityAllManagedUsers), []byte(`"all_managed_users"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.us.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("InvitabilityLevel.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InvitabilityLevel.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberViewabilityLevel_String(t *testing.T) {
	tests := []struct {
		name string
		us   *MemberViewabilityLevel
		want string
	}{
		{"nil", nil, "<nil>"},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAdminsOnly), "admins_only"},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAdminsMembers), "admins_and_members"},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAllManagedUsers), "all_managed_users"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.String(); got != tt.want {
				t.Errorf("MemberViewabilityLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemberViewabilityLevel_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		us      *MemberViewabilityLevel
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte("null"), false},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAdminsOnly), []byte(`"admins_only"`), false},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAdminsMembers), []byte(`"admins_and_members"`), false},
		{"normal", setMemberViewabilityLevel(MemberViewabilityAllManagedUsers), []byte(`"all_managed_users"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.us.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MemberViewabilityLevel.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemberViewabilityLevel.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_ResourceType(t *testing.T) {
	tests := []struct {
		name   string
		target *Group
		want   BoxResourceType
	}{
		// TODO: Add test cases.
		{"normal", nil, GroupResource},
		{"normal", &Group{}, GroupResource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.target
			if got := g.ResourceType(); got != tt.want {
				t.Errorf("Group.ResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGroup(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		api *ApiConn
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		{"normal", args{api: apiConn}, &Group{apiInfo: &apiInfo{api: apiConn}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGroup(tt.args.api); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_GetGroupReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		groupId string
		fields  []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		// TODO: Add test cases.
		{
			name: "normal/fields=nil",
			args: args{groupId: "10001", fields: nil},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10001",
				Method:             GET,
				headers:            http.Header{},
				body:               nil,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
			},
		},
		{
			name: "normal/fields",
			args: args{groupId: "10002", fields: []string{"type", "id"}},
			want: &Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10002?fields=type,id",
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
			t.Helper()

			g := NewGroup(apiConn)
			got := g.GetGroupReq(tt.args.groupId, tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got)
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("differ:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestGroup_GetGroup(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/groups") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/groups")
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
			groupId := strings.TrimPrefix(r.URL.Path, "/2.0/groups/")

			switch groupId {
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
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/groups/group_normal.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildGroupOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		groupId string
		fields  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Group
		wantErr bool
		errType interface{}
	}{
		// TODO: Add test cases.
		{"normal", args{"10001", nil}, normal, false, nil},
		{"http error/404", args{"404", FolderAllFields}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", nil}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", nil}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			g := NewGroup(apiConn)
			got, err := g.GetGroup(tt.args.groupId, tt.args.fields)

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
			opts := diffCompOptions(*got, apiInfo{})
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
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

func TestGroup_CreateGroupReq(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
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
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.CreateGroupReq(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.CreateGroupReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_CreateGroup(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			got, err := g.CreateGroup(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.CreateGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetName(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetProvenance(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		provenance string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetProvenance(tt.args.provenance); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetProvenance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetExternalSyncIdentifier(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		externalSyncIdentifier string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetExternalSyncIdentifier(tt.args.externalSyncIdentifier); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetExternalSyncIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetDescription(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		description string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetDescription(tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetInvitabilityLevel(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		invitabilityLevel InvitabilityLevel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetInvitabilityLevel(tt.args.invitabilityLevel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetInvitabilityLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_SetMemberViewabiityLevel(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		memberViewabilityLevel MemberViewabilityLevel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.SetMemberViewabiityLevel(tt.args.memberViewabilityLevel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.SetMemberViewabiityLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_UpdateGroupReq(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		groupId string
		fields  []string
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
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.UpdateGroupReq(tt.args.groupId, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.UpdateGroupReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_UpdateGroup(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		groupId string
		fields  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Group
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			got, err := g.UpdateGroup(tt.args.groupId, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.UpdateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.UpdateGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_DeleteGroupReq(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		groupId string
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
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.DeleteGroupReq(tt.args.groupId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.DeleteGroupReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_DeleteGroup(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		groupId string
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
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if err := g.DeleteGroup(tt.args.groupId); (err != nil) != tt.wantErr {
				t.Errorf("Group.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroup_GetEnterpriseGroupsReq(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		name   string
		offset int32
		limit  int32
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
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			if got := g.GetEnterpriseGroupsReq(tt.args.name, tt.args.offset, tt.args.limit, tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.GetEnterpriseGroupsReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroup_GetEnterpriseGroups(t *testing.T) {
	type fields struct {
		UserGroupMini          UserGroupMini
		apiInfo                *apiInfo
		CreatedAt              *time.Time
		ModifiedAt             *time.Time
		Provenance             *string
		ExternalSyncIdentifier *string
		Description            *string
		InvitabilityLevel      *InvitabilityLevel
		MemberViewabilityLevel *MemberViewabilityLevel
	}
	type args struct {
		name   string
		offset int32
		limit  int32
		fields []string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantOutGroups     []*Group
		wantOutOffset     int
		wantOutLimit      int
		wantOutTotalCount int
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Group{
				UserGroupMini:          tt.fields.UserGroupMini,
				apiInfo:                tt.fields.apiInfo,
				CreatedAt:              tt.fields.CreatedAt,
				ModifiedAt:             tt.fields.ModifiedAt,
				Provenance:             tt.fields.Provenance,
				ExternalSyncIdentifier: tt.fields.ExternalSyncIdentifier,
				Description:            tt.fields.Description,
				InvitabilityLevel:      tt.fields.InvitabilityLevel,
				MemberViewabilityLevel: tt.fields.MemberViewabilityLevel,
			}
			gotOutGroups, gotOutOffset, gotOutLimit, gotOutTotalCount, err := g.GetEnterpriseGroups(tt.args.name, tt.args.offset, tt.args.limit, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.GetEnterpriseGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutGroups, tt.wantOutGroups) {
				t.Errorf("Group.GetEnterpriseGroups() gotOutGroups = %v, want %v", gotOutGroups, tt.wantOutGroups)
			}
			if gotOutOffset != tt.wantOutOffset {
				t.Errorf("Group.GetEnterpriseGroups() gotOutOffset = %v, want %v", gotOutOffset, tt.wantOutOffset)
			}
			if gotOutLimit != tt.wantOutLimit {
				t.Errorf("Group.GetEnterpriseGroups() gotOutLimit = %v, want %v", gotOutLimit, tt.wantOutLimit)
			}
			if gotOutTotalCount != tt.wantOutTotalCount {
				t.Errorf("Group.GetEnterpriseGroups() gotOutTotalCount = %v, want %v", gotOutTotalCount, tt.wantOutTotalCount)
			}
		})
	}
}
