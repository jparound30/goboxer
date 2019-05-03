package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func setMembershipRole(role MembershipRole) *MembershipRole {
	return &role
}

func TestMembership_Unmarshal(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2013-05-16T15:27:57-07:00")
	modifiedAt, _ := time.Parse(time.RFC3339, "2013-05-16T15:27:57-07:00")

	typ := "group_membership"
	id := "1560354"
	usertype := TYPE_USER
	userid := "13130406"
	username := "Alison Wonderland"
	userlogin := "alice@gmail.com"
	grouptype := TYPE_GROUP
	groupid := "119720"
	groupname := "family"
	role := MembershipRoleAdmin
	tr := true
	fa := false
	configPerm := ConfigurablePermissions{
		CanRunReports:     &fa,
		CanInstantLogin:   &tr,
		CanCreateAccounts: &fa,
		CanEditAccounts:   &tr,
	}

	user := UserGroupMini{
		Type:  &usertype,
		ID:    &userid,
		Name:  &username,
		Login: &userlogin,
	}
	group := UserGroupMini{
		Type: &grouptype,
		ID:   &groupid,
		Name: &groupname,
	}
	tests := []struct {
		name     string
		jsonFile string
		want     Membership
	}{
		{
			name:     "normal",
			jsonFile: "testdata/membership/membership_json.json",
			want: Membership{
				nil,
				&typ,
				&id,
				&user,
				&group,
				&role,
				&createdAt,
				&modifiedAt,
				&configPerm,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := ioutil.ReadFile(tt.jsonFile)
			membership := Membership{}
			err := json.Unmarshal(b, &membership)
			if err != nil {
				t.Errorf("Membership Unmarshal err %v", err)
			}
			if !reflect.DeepEqual(&tt.want, &membership) {
				t.Errorf("Membership Marshal/Unmarshal = %+v, want %+v", membership, tt.want)
			}
		})
	}
}

func TestMembershipRole_String(t *testing.T) {
	tests := []struct {
		name string
		us   *MembershipRole
		want string
	}{
		{"nil", nil, "<nil>"},
		{"normal/admin", setMembershipRole(MembershipRoleAdmin), `admin`},
		{"normal/member", setMembershipRole(MembershipRoleMember), `member`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.String(); got != tt.want {
				t.Errorf("MembershipRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembershipRole_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		us      *MembershipRole
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"normal/admin", setMembershipRole(MembershipRoleAdmin), []byte(`"admin"`), false},
		{"normal/member", setMembershipRole(MembershipRoleMember), []byte(`"member"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.us.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MembershipRole.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MembershipRole.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembership_ResourceType(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	tests := []struct {
		name string
		want BoxResourceType
	}{
		{"normal", MembershipResource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMembership(apiConn)
			if got := m.ResourceType(); got != tt.want {
				t.Errorf("Membership.ResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMembership_GetMembershipReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		membershipId string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal", args{"10001"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships/10001",
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
			m := NewMembership(apiConn)
			got := m.GetMembershipReq(tt.args.membershipId)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}

		})
	}
}

func TestMembership_GetMembership(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/group_memberships") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/group_memberships")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/membership_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildMembershipOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	type args struct {
		membershipId string
	}
	tests := []struct {
		name    string
		args    args
		want    *Membership
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001"}, normal, false, nil},
		{"http error/404", args{"404"}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999"}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999"}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			m := NewMembership(apiConn)
			got, err := m.GetMembership(tt.args.membershipId)

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

func buildMembershipOfGetInfoNormalJson() *Membership {
	var normal Membership

	normal.Type = setStringPtr("group_membership")
	normal.ID = setStringPtr("1560354")

	normal.User = &UserGroupMini{
		Type:  setUserType(TYPE_USER),
		ID:    setStringPtr("13130406"),
		Name:  setStringPtr("Alison Wonderland"),
		Login: setStringPtr("alice@gmail.com"),
	}
	normal.Group = &UserGroupMini{
		Type: setUserType(TYPE_GROUP),
		ID:   setStringPtr("119720"),
		Name: setStringPtr("family"),
	}
	normal.Role = setMembershipRole(MembershipRoleAdmin)
	normal.ConfigurablePermissions = &ConfigurablePermissions{
		CanRunReports:     setBool(false),
		CanInstantLogin:   setBool(true),
		CanCreateAccounts: setBool(false),
		CanEditAccounts:   setBool(true),
	}
	normal.CreatedAt = setTime("2013-05-16T15:27:57-07:00")
	normal.ModifiedAt = setTime("2013-05-16T15:27:57-07:00")

	return &normal
}

func TestMembership_CreateMembershipReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type fields struct {
		apiInfo                 *apiInfo
		Type                    *string
		ID                      *string
		User                    *UserGroupMini
		Group                   *UserGroupMini
		Role                    *MembershipRole
		CreatedAt               *time.Time
		ModifiedAt              *time.Time
		ConfigurablePermissions *ConfigurablePermissions
	}
	tests := []struct {
		name   string
		fields fields
		want   *Request
	}{
		{"normal",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				ID:      setStringPtr("i10001"),
				User: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("ui10001"),
					Name:  setStringPtr("un10001"),
					Login: setStringPtr("ul10001"),
				},
				Group: &UserGroupMini{
					Type: setUserType(TYPE_GROUP),
					ID:   setStringPtr("gi10001"),
					Name: setStringPtr("gn10001"),
				},
				Role:       setMembershipRole(MembershipRoleAdmin),
				CreatedAt:  setTime("2013-05-16T15:27:57-07:00"),
				ModifiedAt: setTime("2013-06-16T15:27:57-07:00"),
				ConfigurablePermissions: &ConfigurablePermissions{
					CanRunReports:     setBool(true),
					CanInstantLogin:   setBool(false),
					CanCreateAccounts: setBool(true),
					CanEditAccounts:   setBool(false),
				},
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"user": {
		"id": "ui10001"
	},
	"group": {
		"id": "gi10001"
	},
	"role": "admin",
	"configurable_permissions": {
		"can_run_reports": true,
		"can_instant_login": false,
		"can_create_accounts": true,
		"can_edit_accounts": false
	}
}
`),
			},
		},
		{"user/group",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				User: &UserGroupMini{
					ID: setStringPtr("ui10002"),
				},
				Group: &UserGroupMini{
					ID: setStringPtr("gi10002"),
				},
				Role:                    nil,
				ConfigurablePermissions: nil,
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"user": {
		"id": "ui10002"
	},
	"group": {
		"id": "gi10002"
	}
}
`),
			},
		},
		{"role",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				User: &UserGroupMini{
					ID: setStringPtr("ui10003"),
				},
				Group: &UserGroupMini{
					ID: setStringPtr("gi10003"),
				},
				Role: setMembershipRole(MembershipRoleMember),
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"user": {
		"id": "ui10003"
	},
	"group": {
		"id": "gi10003"
	},
	"role": "member"
}
`),
			},
		},
		{"configuration",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				User: &UserGroupMini{
					ID: setStringPtr("ui10004"),
				},
				Group: &UserGroupMini{
					Type: setUserType(TYPE_GROUP),
					ID:   setStringPtr("gi10004"),
				},
				ConfigurablePermissions: &ConfigurablePermissions{
					CanRunReports: setBool(false),
				},
			},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"user": {
		"id": "ui10004"
	},
	"group": {
		"id": "gi10004"
	},
	"configurable_permissions": {
		"can_run_reports": false
	}
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Membership{
				apiInfo:                 tt.fields.apiInfo,
				Type:                    tt.fields.Type,
				ID:                      tt.fields.ID,
				User:                    tt.fields.User,
				Group:                   tt.fields.Group,
				Role:                    tt.fields.Role,
				CreatedAt:               tt.fields.CreatedAt,
				ModifiedAt:              tt.fields.ModifiedAt,
				ConfigurablePermissions: tt.fields.ConfigurablePermissions,
			}

			got := m.CreateMembershipReq()

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

func TestMembership_CreateMembership(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/group_memberships") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/group_memberships")
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
			userId := js["user"].(map[string]interface{})["id"]

			switch userId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/membership_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildMembershipOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	m1 := buildMembershipOfGetInfoNormalJson()
	m1.apiInfo = &apiInfo{api: apiConn}
	m1.SetUser("u10001")
	m1.SetGroup("g10001")
	m1.SetRole(MembershipRoleMember)
	m1.SetConfigurablePermissions(true, false, true, false)

	m2 := buildMembershipOfGetInfoNormalJson()
	m2.apiInfo = &apiInfo{api: apiConn}
	m2.SetUser("404")

	m3 := buildMembershipOfGetInfoNormalJson()
	m3.apiInfo = &apiInfo{api: apiConn}
	m3.SetUser("999")

	m4 := buildMembershipOfGetInfoNormalJson()
	m4.apiInfo = &apiInfo{api: apiConn}
	m4.SetUser("999")

	tests := []struct {
		name    string
		target  *Membership
		want    *Membership
		wantErr bool
		errType interface{}
	}{
		{"normal",
			m1,
			normal,
			false,
			nil,
		},
		{"http error/404",
			m2,
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			m3,
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			m4,
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

			m := tt.target
			got, err := m.CreateMembership()

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

func TestMembership_UpdateMembershipReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type fields struct {
		apiInfo                 *apiInfo
		Type                    *string
		ID                      *string
		User                    *UserGroupMini
		Group                   *UserGroupMini
		Role                    *MembershipRole
		CreatedAt               *time.Time
		ModifiedAt              *time.Time
		ConfigurablePermissions *ConfigurablePermissions
	}
	type args struct {
		membershipId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{"normal",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				ID:      setStringPtr("i10001"),
				User: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("ui10001"),
					Name:  setStringPtr("un10001"),
					Login: setStringPtr("ul10001"),
				},
				Group: &UserGroupMini{
					Type: setUserType(TYPE_GROUP),
					ID:   setStringPtr("gi10001"),
					Name: setStringPtr("gn10001"),
				},
				Role:       setMembershipRole(MembershipRoleAdmin),
				CreatedAt:  setTime("2013-05-16T15:27:57-07:00"),
				ModifiedAt: setTime("2013-06-16T15:27:57-07:00"),
				ConfigurablePermissions: &ConfigurablePermissions{
					CanRunReports:     setBool(true),
					CanInstantLogin:   setBool(false),
					CanCreateAccounts: setBool(true),
					CanEditAccounts:   setBool(false),
				},
			},
			args{"10001"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"role": "admin",
	"configurable_permissions": {
		"can_run_reports": true,
		"can_instant_login": false,
		"can_create_accounts": true,
		"can_edit_accounts": false
	}
}
`),
			},
		},
		{"role",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				User: &UserGroupMini{
					ID: setStringPtr("ui10003"),
				},
				Group: &UserGroupMini{
					ID: setStringPtr("gi10003"),
				},
				Role: setMembershipRole(MembershipRoleMember),
			},
			args{"10002"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships/10002",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"role": "member"
}
`),
			},
		},
		{"configuration",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				User: &UserGroupMini{
					ID: setStringPtr("ui10004"),
				},
				Group: &UserGroupMini{
					Type: setUserType(TYPE_GROUP),
					ID:   setStringPtr("gi10004"),
				},
				ConfigurablePermissions: &ConfigurablePermissions{
					CanRunReports: setBool(false),
				},
			},
			args{"10003"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships/10003",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"configurable_permissions": {
		"can_run_reports": false
	}
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Membership{
				apiInfo:                 tt.fields.apiInfo,
				Type:                    tt.fields.Type,
				ID:                      tt.fields.ID,
				User:                    tt.fields.User,
				Group:                   tt.fields.Group,
				Role:                    tt.fields.Role,
				CreatedAt:               tt.fields.CreatedAt,
				ModifiedAt:              tt.fields.ModifiedAt,
				ConfigurablePermissions: tt.fields.ConfigurablePermissions,
			}

			got := m.UpdateMembershipReq(tt.args.membershipId)

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

func TestMembership_UpdateMembership(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/group_memberships") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/group_memberships")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/membership_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildMembershipOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	m1 := buildMembershipOfGetInfoNormalJson()
	m1.apiInfo = &apiInfo{api: apiConn}
	m1.SetUser("u10001")
	m1.SetGroup("g10001")
	m1.SetRole(MembershipRoleMember)
	m1.SetConfigurablePermissions(true, false, true, false)

	m2 := buildMembershipOfGetInfoNormalJson()
	m2.apiInfo = &apiInfo{api: apiConn}
	m2.SetUser("404")

	m3 := buildMembershipOfGetInfoNormalJson()
	m3.apiInfo = &apiInfo{api: apiConn}
	m3.SetUser("999")

	m4 := buildMembershipOfGetInfoNormalJson()
	m4.apiInfo = &apiInfo{api: apiConn}
	m4.SetUser("999")

	type args struct {
		membershipId string
	}
	tests := []struct {
		name    string
		target  *Membership
		args    args
		want    *Membership
		wantErr bool
		errType interface{}
	}{
		{"normal",
			m1,
			args{"10001"},
			normal,
			false,
			nil,
		},
		{"http error/404",
			m2,
			args{"404"},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"returned invalid json/999",
			m3,
			args{"999"},
			nil,
			true,
			&ApiOtherError{}},
		{"senderror",
			m4,
			args{"999"},
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

			m := tt.target
			got, err := m.UpdateMembership(tt.args.membershipId)

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

func TestMembership_DeleteMembershipReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type fields struct {
		apiInfo                 *apiInfo
		Type                    *string
		ID                      *string
		User                    *UserGroupMini
		Group                   *UserGroupMini
		Role                    *MembershipRole
		CreatedAt               *time.Time
		ModifiedAt              *time.Time
		ConfigurablePermissions *ConfigurablePermissions
	}
	type args struct {
		membershipId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Request
	}{
		{"normal",
			fields{
				apiInfo: &apiInfo{api: apiConn},
				Type:    setStringPtr("group_membership"),
				ID:      setStringPtr("i10001"),
				User: &UserGroupMini{
					Type:  setUserType(TYPE_USER),
					ID:    setStringPtr("ui10001"),
					Name:  setStringPtr("un10001"),
					Login: setStringPtr("ul10001"),
				},
				Group: &UserGroupMini{
					Type: setUserType(TYPE_GROUP),
					ID:   setStringPtr("gi10001"),
					Name: setStringPtr("gn10001"),
				},
				Role:       setMembershipRole(MembershipRoleAdmin),
				CreatedAt:  setTime("2013-05-16T15:27:57-07:00"),
				ModifiedAt: setTime("2013-06-16T15:27:57-07:00"),
				ConfigurablePermissions: &ConfigurablePermissions{
					CanRunReports:     setBool(true),
					CanInstantLogin:   setBool(false),
					CanCreateAccounts: setBool(true),
					CanEditAccounts:   setBool(false),
				},
			},
			args{"10001"},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/group_memberships/10001",
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
			m := &Membership{
				apiInfo:                 tt.fields.apiInfo,
				Type:                    tt.fields.Type,
				ID:                      tt.fields.ID,
				User:                    tt.fields.User,
				Group:                   tt.fields.Group,
				Role:                    tt.fields.Role,
				CreatedAt:               tt.fields.CreatedAt,
				ModifiedAt:              tt.fields.ModifiedAt,
				ConfigurablePermissions: tt.fields.ConfigurablePermissions,
			}

			got := m.DeleteMembershipReq(tt.args.membershipId)

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

func TestMembership_DeleteMembership(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/group_memberships") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/group_memberships")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(204)
				resp, _ := ioutil.ReadFile("testdata/membership/membership_json.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := buildMembershipOfGetInfoNormalJson()
	normal.apiInfo = &apiInfo{api: apiConn}

	m1 := buildMembershipOfGetInfoNormalJson()
	m1.apiInfo = &apiInfo{api: apiConn}
	m1.SetUser("u10001")
	m1.SetGroup("g10001")
	m1.SetRole(MembershipRoleMember)
	m1.SetConfigurablePermissions(true, false, true, false)

	m2 := buildMembershipOfGetInfoNormalJson()
	m2.apiInfo = &apiInfo{api: apiConn}
	m2.SetUser("404")

	m4 := buildMembershipOfGetInfoNormalJson()
	m4.apiInfo = &apiInfo{api: apiConn}
	m4.SetUser("999")

	type args struct {
		membershipId string
	}
	tests := []struct {
		name    string
		target  *Membership
		args    args
		want    *Membership
		wantErr bool
		errType interface{}
	}{
		{"normal",
			m1,
			args{"10001"},
			normal,
			false,
			nil,
		},
		{"http error/404",
			m2,
			args{"404"},
			nil,
			true,
			&ApiStatusError{Status: 404},
		},
		{"senderror",
			m4,
			args{"999"},
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

			m := tt.target
			err := m.DeleteMembership(tt.args.membershipId)

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

func TestMembership_GetMembershipForGroupReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		groupId string
		offset  int32
		limit   int32
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", 0, 1000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10001/memberships?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10002", 0, 2000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10002/memberships?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10003", 900, 100},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10003/memberships?offset=900&limit=100",
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
			m := NewMembership(apiConn)
			got := m.GetMembershipForGroupReq(tt.args.groupId, tt.args.offset, tt.args.limit)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestMembership_GetMembershipForGroup(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/groups/") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/groups/")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/membership_for_group.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		groupId string
		offset  int32
		limit   int32
	}
	tests := []struct {
		name    string
		args    args
		want    []*Membership
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", 0, 1000},
			[]*Membership{
				{
					apiInfo: &apiInfo{api: apiConn},
					Type:    setStringPtr("group_membership"),
					ID:      setStringPtr("1560354"),
					User: &UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("13130906"),
						Name:  setStringPtr("Alice"),
						Login: setStringPtr("alice@gmail.com"),
					},
					Group: &UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("119720"),
						Name: setStringPtr("family"),
					},
					Role: setMembershipRole(MembershipRoleMember),
				},
				{
					apiInfo: &apiInfo{api: apiConn},
					Type:    setStringPtr("group_membership"),
					ID:      setStringPtr("1560356"),
					User: &UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("192633962"),
						Name:  setStringPtr("rabbit"),
						Login: setStringPtr("rabbit@gmail.com"),
					},
					Group: &UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("119720"),
						Name: setStringPtr("family"),
					},
					Role: setMembershipRole(MembershipRoleMember),
				},
			}, false, nil},
		{"http error/404", args{"404", 0, 1000}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			m := NewMembership(apiConn)
			got, gotOffset, gotLimit, gotTotalCount, err := m.GetMembershipForGroup(tt.args.groupId, tt.args.offset, tt.args.limit)

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
			if gotTotalCount != 2 {
				t.Errorf("totalCount incorrect")
			}
			if gotOffset != 0 {
				t.Errorf("offset incorrect")
			}
			if gotLimit != 100 {
				t.Errorf("limi incorrect")
			}
			opts := diffCompOptions(Membership{}, apiInfo{})
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			for _, v := range got {
				if v.apiInfo == nil {
					t.Errorf("not exists `apiInfo` field\n")
				}
				return
			}
		})
	}
}

func TestMembership_GetMembershipForUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		userId string
		offset int32
		limit  int32
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", 0, 1000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10001/memberships?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10002", 0, 2000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10002/memberships?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10003", 900, 100},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10003/memberships?offset=900&limit=100",
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
			m := NewMembership(apiConn)
			got := m.GetMembershipForUserReq(tt.args.userId, tt.args.offset, tt.args.limit)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestMembership_GetMembershipForUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users/") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users/")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/membership_for_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		userId string
		offset int32
		limit  int32
	}
	tests := []struct {
		name    string
		args    args
		want    []*Membership
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", 0, 1000},
			[]*Membership{
				{
					apiInfo: &apiInfo{api: apiConn},
					Type:    setStringPtr("group_membership"),
					ID:      setStringPtr("1560354"),
					User: &UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("13130406"),
						Name:  setStringPtr("Alison Wonderland"),
						Login: setStringPtr("alice@gmail.com"),
					},
					Group: &UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("119720"),
						Name: setStringPtr("family"),
					},
					Role: setMembershipRole(MembershipRoleMember),
				},
			}, false, nil},
		{"http error/404", args{"404", 0, 1000}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			m := NewMembership(apiConn)
			got, gotOffset, gotLimit, gotTotalCount, err := m.GetMembershipForUser(tt.args.userId, tt.args.offset, tt.args.limit)

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
			if gotTotalCount != 1 {
				t.Errorf("totalCount incorrect")
			}
			if gotOffset != 0 {
				t.Errorf("offset incorrect")
			}
			if gotLimit != 100 {
				t.Errorf("limi incorrect")
			}
			opts := diffCompOptions(Membership{}, apiInfo{})
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			for _, v := range got {
				if v.apiInfo == nil {
					t.Errorf("not exists `apiInfo` field\n")
				}
				return
			}
		})
	}
}

func TestMembership_GetCollaborationsForGroupReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		groupID string
		offset  int32
		limit   int32
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", 0, 1000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10001/collaborations?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10002", 0, 2000},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10002/collaborations?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10003", 900, 100},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/groups/10003/collaborations?offset=900&limit=100",
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
			m := NewMembership(apiConn)
			got := m.GetCollaborationsForGroupReq(tt.args.groupID, tt.args.offset, tt.args.limit)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestMembership_GetCollaborationsForGroup(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/groups/") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/groups/")
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
			membershipId := strings.Split(r.URL.Path, "/")[3]

			switch membershipId {
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
				resp, _ := ioutil.ReadFile("testdata/membership/collaboration_for_group.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		groupId string
		offset  int32
		limit   int32
	}
	tests := []struct {
		name    string
		args    args
		want    []*Collaboration
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", 0, 1000},
			[]*Collaboration{
				{
					apiInfo: &apiInfo{api: apiConn},
					Type:    setStringPtr("collaboration"),
					ID:      setStringPtr("52123184"),
					CreatedBy: &UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("13130406"),
						Name:  setStringPtr("Eddard Stark"),
						Login: setStringPtr("ned@winterfell.com"),
					},
					CreatedAt:  setTime("2013-11-14T16:16:20-08:00"),
					ModifiedAt: setTime("2013-11-14T16:16:20-08:00"),
					ExpiresAt:  nil,
					Status:     setCollaborationStatus(COLLABORATION_STATUS_ACCEPTED),
					AccessibleBy: &UserGroupMini{
						Type: setUserType(TYPE_GROUP),
						ID:   setStringPtr("160018"),
						Name: setStringPtr("Hand of the King inner counsel"),
					},
					Role:           setRole(VIEWER),
					AcknowledgedAt: setTime("2013-11-14T16:16:20-08:00"),
					Item: &ItemMini{
						Type:       setItemTypePtr(TYPE_FOLDER),
						ID:         setStringPtr("541014843"),
						SequenceId: setStringPtr("0"),
						ETag:       setStringPtr("0"),
						Name:       setStringPtr("People killed by Ice"),
					},
				},
			}, false, nil},
		{"http error/404", args{"404", 0, 1000}, nil, true, &ApiStatusError{Status: 404}},
		{"returned invalid json/999", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
		{"senderror", args{"999", 0, 1000}, nil, true, &ApiOtherError{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			if tt.name == "senderror" {
				apiConn.Expires = 0
			} else {
				apiConn.Expires = 6000
			}

			m := NewMembership(apiConn)
			got, gotOffset, gotLimit, gotTotalCount, err := m.GetCollaborationsForGroup(tt.args.groupId, tt.args.offset, tt.args.limit)

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
			if gotTotalCount != 1 {
				t.Errorf("totalCount incorrect")
			}
			if gotOffset != 0 {
				t.Errorf("offset incorrect")
			}
			if gotLimit != 100 {
				t.Errorf("limi incorrect")
			}
			opts := diffCompOptions(Collaboration{}, apiInfo{})
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			for _, v := range got {
				if v.apiInfo == nil {
					t.Errorf("not exists `apiInfo` field\n")
				}
				return
			}
		})
	}
}
