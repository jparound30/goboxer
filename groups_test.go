package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	url := "https://example.com"
	apiConn := commonInit(url)

	g1 := buildGroupOfGetInfoNormalJson()
	g1.SetName("name1")
	g1.apiInfo = &apiInfo{api: apiConn}

	g2 := buildGroupOfGetInfoNormalJson()
	g2.SetProvenance("provenance1")
	g2.apiInfo = &apiInfo{api: apiConn}

	g3 := buildGroupOfGetInfoNormalJson()
	g3.SetExternalSyncIdentifier("external_id_1")
	g3.apiInfo = &apiInfo{api: apiConn}

	g4 := buildGroupOfGetInfoNormalJson()
	g4.SetDescription("description1")
	g4.apiInfo = &apiInfo{api: apiConn}

	g5 := buildGroupOfGetInfoNormalJson()
	g5.SetInvitabilityLevel(InvitabilityAllManagedUsers)
	g5.apiInfo = &apiInfo{api: apiConn}

	g6 := buildGroupOfGetInfoNormalJson()
	g6.SetMemberViewabiityLevel(MemberViewabilityAdminsMembers)
	g6.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		target *Group
		args   args
		want   *Request
	}{
		{"normal/name",
			g1,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "name1"
}
`),
			},
		},
		{"normal/provenance",
			g2,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"provenance": "provenance1"
}
`),
			},
		},
		{"normal/external sync identifier",
			g3,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"external_sync_identifier": "external_id_1"
}
`),
			},
		},
		{"normal/description",
			g4,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"description": "description1"
}
`),
			},
		},
		{"normal/invitability",
			g5,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"invitability_level": "all_managed_users"
}
`),
			},
		},
		{"normal/member viewability",
			g6,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"member_viewability_level": "admins_and_members"
}
`),
			},
		},
		{"normal/fields",
			g6,
			args{[]string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups?fields=type,id",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"member_viewability_level": "admins_and_members"
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.target
			got := g.CreateGroupReq(tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
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

func TestGroup_CreateGroup(t *testing.T) {
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
			groupName := js["name"]

			switch groupName {
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

	g1 := NewGroup(apiConn)
	g1.SetName("10001")

	g2 := NewGroup(apiConn)
	g2.SetName("404")

	g3 := NewGroup(apiConn)
	g3.SetName("999")

	g4 := NewGroup(apiConn)
	g4.SetName("999")

	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		target  *Group
		args    args
		want    *Group
		wantErr bool
		errType interface{}
	}{
		{"normal",
			g1,
			args{
				[]string{"type", "id"},
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			g2,
			args{
				[]string{"type", "id"},
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			g3,
			args{
				[]string{"type", "id"},
			},
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			g4,
			args{
				[]string{"type", "id"},
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

			g := tt.target
			got, err := g.CreateGroup(tt.args.fields)

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
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
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

func TestGroup_UpdateGroupReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	g1 := buildGroupOfGetInfoNormalJson()
	g1.SetName("name1")
	g1.apiInfo = &apiInfo{api: apiConn}

	g2 := buildGroupOfGetInfoNormalJson()
	g2.SetProvenance("provenance1")
	g2.apiInfo = &apiInfo{api: apiConn}

	g3 := buildGroupOfGetInfoNormalJson()
	g3.SetExternalSyncIdentifier("external_id_1")
	g3.apiInfo = &apiInfo{api: apiConn}

	g4 := buildGroupOfGetInfoNormalJson()
	g4.SetDescription("description1")
	g4.apiInfo = &apiInfo{api: apiConn}

	g5 := buildGroupOfGetInfoNormalJson()
	g5.SetInvitabilityLevel(InvitabilityAllManagedUsers)
	g5.apiInfo = &apiInfo{api: apiConn}

	g6 := buildGroupOfGetInfoNormalJson()
	g6.SetMemberViewabiityLevel(MemberViewabilityAdminsMembers)
	g6.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		groupId string
		fields  []string
	}
	tests := []struct {
		name   string
		target *Group
		args   args
		want   *Request
	}{
		{"normal/name",
			g1,
			args{"10001", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "name1"
}
`),
			},
		},
		{"normal/provenance",
			g2,
			args{"10002", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10002",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"provenance": "provenance1"
}
`),
			},
		},
		{"normal/external sync identifier",
			g3,
			args{"10003", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10003",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"external_sync_identifier": "external_id_1"
}
`),
			},
		},
		{"normal/description",
			g4,
			args{"10004", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10004",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"description": "description1"
}
`),
			},
		},
		{"normal/invitability",
			g5,
			args{"10005", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10005",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"invitability_level": "all_managed_users"
}
`),
			},
		},
		{"normal/member viewability",
			g6,
			args{"10006", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10006",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"member_viewability_level": "admins_and_members"
}
`),
			},
		},
		{"normal/fields",
			g6,
			args{"10007", []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10007?fields=type,id",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"member_viewability_level": "admins_and_members"
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.target
			got := g.UpdateGroupReq(tt.args.groupId, tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
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

func TestGroup_UpdateGroup(t *testing.T) {
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
			if r.Method != http.MethodPut {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			groupId := strings.Split(r.URL.Path, "/")[3]

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

	g1 := NewGroup(apiConn)
	g1.SetName("10001")

	g2 := NewGroup(apiConn)
	g2.SetName("404")

	g3 := NewGroup(apiConn)
	g3.SetName("999")

	g4 := NewGroup(apiConn)
	g4.SetName("999")

	type args struct {
		groupId string
		fields  []string
	}
	tests := []struct {
		name    string
		target  *Group
		args    args
		want    *Group
		wantErr bool
		errType interface{}
	}{
		{"normal",
			g1,
			args{
				"10001", []string{"type", "id"},
			},
			normal,
			false,
			nil,
		},
		{"http error/404",
			g2,
			args{
				"404", []string{"type", "id"},
			},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			g3,
			args{
				"999", []string{"type", "id"},
			},
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			g4,
			args{
				"999", []string{"type", "id"},
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

			g := tt.target
			got, err := g.UpdateGroup(tt.args.groupId, tt.args.fields)

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
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
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

func TestGroup_DeleteGroupReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	g1 := buildGroupOfGetInfoNormalJson()
	g1.SetName("name1")
	g1.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		groupId string
	}
	tests := []struct {
		name   string
		target *Group
		args   args
		want   *Request
	}{
		{"normal", g1, args{"10001"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10001",
				Method:             DELETE,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			g := tt.target
			got := g.DeleteGroupReq(tt.args.groupId)

			// If normal response
			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestGroup_DeleteGroup(t *testing.T) {
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
			if r.Method != http.MethodDelete {
				t.Fatalf("invalid http method")
			}
			// Header check
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("not exists access token")
			}
			// ok, return some response
			groupId := strings.Split(r.URL.Path, "/")[3]

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

	normal := buildGroupOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		groupId string
	}
	tests := []struct {
		name    string
		args    args
		want    *Group
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

			g := NewGroup(apiConn)
			err := g.DeleteGroup(tt.args.groupId)

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

func TestGroup_GetEnterpriseGroupsReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		name   string
		offset int32
		limit  int32
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		// TODO: Add test cases.
		{"no name",
			args{
				"",
				0,
				2000,
				nil,
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"name",
			args{
				"CONTAINS",
				0,
				100,
				nil,
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups?offset=0&limit=100&name=CONTAINS",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"name escape",
			args{
				"CON TAINS",
				900,
				100,
				nil,
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups?offset=900&limit=100&name=CON+TAINS",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"fields",
			args{
				"あいうえお",
				900,
				100,
				[]string{"type", "id"},
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups?offset=900&limit=100&name=%E3%81%82%E3%81%84%E3%81%86%E3%81%88%E3%81%8A&fields=type,id",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			g := NewGroup(apiConn)
			got := g.GetEnterpriseGroupsReq(tt.args.name, tt.args.offset, tt.args.limit, tt.args.fields)

			// If normal response
			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestGroup_GetEnterpriseGroups(t *testing.T) {
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
			name := r.URL.Query().Get("name")

			switch name {
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
				resp, _ := ioutil.ReadFile("testdata/groups/get_enterprise_normal.json")
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
		name   string
		offset int32
		limit  int32
		fields []string
	}
	tests := []struct {
		name              string
		args              args
		wantOutGroups     []*Group
		wantOutOffset     int
		wantOutLimit      int
		wantOutTotalCount int
		wantErr           bool
		errType           interface{}
	}{
		// TODO: Add test cases.
		{
			"normal",
			args{"10001", 0, 100, nil},
			[]*Group{
				{
					UserGroupMini: UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("1786931"),
						Name: setStringPtr("friends"),
					},
				},
				{
					UserGroupMini: UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("1786932"),
						Name: setStringPtr("あいうえお"),
					},
				},
			},
			0,
			100,
			2,
			false,
			nil,
		},
		{
			"http error/404",
			args{"404", 0, 100, nil},
			nil,
			0,
			100,
			2,
			true,
			&ApiStatusError{Status: 404},
		},
		{
			"returned invalid json/999",
			args{"999", 0, 100, nil},
			nil,
			0,
			100,
			2,
			true,
			&ApiOtherError{},
		},
		{
			"senderror",
			args{"999", 0, 100, nil},
			nil,
			0,
			100,
			2,
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

			g := NewGroup(apiConn)
			gotOutGroups, gotOutOffset, gotOutLimit, gotOutTotalCount, err := g.GetEnterpriseGroups(tt.args.name, tt.args.offset, tt.args.limit, tt.args.fields)

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
			opts := diffCompOptions(Group{})
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(gotOutGroups, tt.wantOutGroups, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
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
