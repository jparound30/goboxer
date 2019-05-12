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

func TestUser_Unmarshal(t *testing.T) {
	typ := TYPE_USER
	id := "181216415"
	name := "sean rose"
	login := "sean+awesome@box.com"
	createdAt, _ := time.Parse(time.RFC3339, "2012-05-03T21:39:11-07:00")
	modifiedAt, _ := time.Parse(time.RFC3339, "2012-11-14T11:21:32-08:00")
	role := UserRoleAdmin
	language := "en"
	timezone := "Africa/Bujumbura"
	spaceAmount := int64(11345156112)
	spaceUsed := int64(1237009912)
	maxUploadSize := 2147483648
	trackingCodes := []map[string]string{}
	canSeeManagedUsers := true
	isSyncEnabled := true
	status := UserStatusActive
	jobTItle := ""
	phone := "6509241374"
	address := ""
	avatarUrl := "https://www.box.com/api/avatar/large/181216415"
	isExemptFromDeviceLimits := false
	isExemptFromLoginVerification := false
	enterprise := Enterprise{
		EnterpriseTypeEnterprise,
		"17077211",
		"seanrose enterprise",
	}
	myTags := []string{"important", "needs review"}
	tests := []struct {
		name     string
		jsonfile string
		want     User
	}{
		{
			name:     "normal",
			jsonfile: "testdata/users/user_json.json",
			want: User{
				UserGroupMini{
					Type:  &typ,
					ID:    &id,
					Name:  &name,
					Login: &login,
				},
				nil,
				&createdAt,
				&modifiedAt,
				&language,
				&timezone,
				spaceAmount,
				spaceUsed,
				maxUploadSize,
				&status,
				&jobTItle,
				&phone,
				&address,
				&avatarUrl,
				&role,
				trackingCodes,
				&canSeeManagedUsers,
				&isSyncEnabled,
				nil,
				&isExemptFromDeviceLimits,
				&isExemptFromLoginVerification,
				&enterprise,
				&myTags,
				nil,
				nil,
				nil,
				nil,
				nil,
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := ioutil.ReadFile(tt.jsonfile)
			user := User{}
			err := json.Unmarshal(b, &user)
			if err != nil {
				t.Errorf("User Unmarshal err %v", err)
			}
			if !reflect.DeepEqual(&tt.want, &user) {
				t.Errorf("User Marshal/Unmarshal = %v, want %v", &user, tt.want)
			}
		})
	}
}

func TestUser_UnmarshalMarshal(t *testing.T) {
	typ := TYPE_USER
	id := "181216415"
	name := "sean rose"
	login := "sean+awesome@box.com"
	createdAt, _ := time.Parse(time.RFC3339, "2012-05-03T21:39:11-07:00")
	modifiedAt, _ := time.Parse(time.RFC3339, "2012-11-14T11:21:32-08:00")
	role := UserRoleAdmin
	language := "en"
	timezone := "Africa/Bujumbura"
	spaceAmount := int64(11345156112)
	spaceUsed := int64(1237009912)
	maxUploadSize := 2147483648
	trackingCodes := []map[string]string{{"k1": "v1"}, {"k2": "v2"}}
	canSeeManagedUsers := true
	isSyncEnabled := true
	status := UserStatusActive
	jobTItle := ""
	phone := "6509241374"
	address := ""
	avatarUrl := "https://www.box.com/api/avatar/large/181216415"
	isExemptFromDeviceLimits := false
	isExemptFromLoginVerification := false
	enterprise := Enterprise{
		EnterpriseTypeEnterprise,
		"17077211",
		"seanrose enterprise",
	}
	myTags := []string{"important", "needs review"}
	tests := []struct {
		name string
		want User
	}{
		{
			name: "normal",
			want: User{
				UserGroupMini{
					Type:  &typ,
					ID:    &id,
					Name:  &name,
					Login: &login,
				},
				nil,
				&createdAt,
				&modifiedAt,
				&language,
				&timezone,
				spaceAmount,
				spaceUsed,
				maxUploadSize,
				&status,
				&jobTItle,
				&phone,
				&address,
				&avatarUrl,
				&role,
				trackingCodes,
				&canSeeManagedUsers,
				&isSyncEnabled,
				nil,
				&isExemptFromDeviceLimits,
				&isExemptFromLoginVerification,
				&enterprise,
				&myTags,
				nil,
				nil,
				nil,
				nil,
				nil,
				0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(&tt.want)
			if err != nil {
				t.Errorf("User Marshal err %v", err)
			}
			u := User{}
			err = json.Unmarshal(b, &u)
			if err != nil {
				t.Errorf("User Unmarshal err %v", err)
			}
			if !reflect.DeepEqual(&tt.want, &u) {
				t.Errorf("User Marshal/Unmarshal = %v, want %v", u, tt.want)
			}
		})
	}
}
func TestUserStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		us      *UserStatus
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"normal/active", setUserStatus(UserStatusActive), []byte(`"active"`), false},
		{"normal/inactive", setUserStatus(UserStatusInactive), []byte(`"inactive"`), false},
		{"normal/cannot_delete_edit", setUserStatus(UserStatusCannotDeleteEdit), []byte(`"cannot_delete_edit"`), false},
		{"normal/cannot_delete_edit_upload", setUserStatus(UserStatusCannotDeleteEditUpload), []byte(`"cannot_delete_edit_upload"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.us.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStatus.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUserStatus(status UserStatus) *UserStatus {
	return &status
}

func TestUserRole_String(t *testing.T) {
	tests := []struct {
		name string
		ur   *UserRole
		want string
	}{
		{"nil", nil, "<nil>"},
		{"admin", setUserRole(UserRoleAdmin), "admin"},
		{"coadmin", setUserRole(UserRoleCoAdmin), "coadmin"},
		{"user", setUserRole(UserRoleUser), "user"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ur.String(); got != tt.want {
				t.Errorf("UserRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRole_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		ur      *UserRole
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"normal/admin", setUserRole(UserRoleAdmin), []byte(`"admin"`), false},
		{"normal/coadmin", setUserRole(UserRoleCoAdmin), []byte(`"coadmin"`), false},
		{"normal/user", setUserRole(UserRoleUser), []byte(`"user"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ur.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRole.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRole.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setUserRole(role UserRole) *UserRole {
	return &role
}

func TestEnterpriseType_String(t *testing.T) {
	tests := []struct {
		name string
		et   *EnterpriseType
		want string
	}{
		{"nil", nil, "<nil>"},
		{"enterprise", setEnterpriseType(EnterpriseTypeEnterprise), "enterprise"},
		{"user", setEnterpriseType(EnterpriseTypeUser), "user"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.String(); got != tt.want {
				t.Errorf("EnterpriseType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnterpriseType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		et      *EnterpriseType
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"normal/enterprise", setEnterpriseType(EnterpriseTypeEnterprise), []byte(`"enterprise"`), false},
		{"normal/user", setEnterpriseType(EnterpriseTypeUser), []byte(`"user"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.et.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("EnterpriseType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnterpriseType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setEnterpriseType(enterpriseType EnterpriseType) *EnterpriseType {
	return &enterpriseType
}

func TestUser_GetCurrentUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal", args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/me",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal", args{[]string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/me?fields=type,id",
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
			u := NewUser(apiConn)
			got := u.GetCurrentUserReq(tt.args.fields)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}

		})
	}
}

func TestUser_GetCurrentUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users/me") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users/me")
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
			fields := r.URL.Query().Get("fields")

			switch fields {
			case "500":
				w.WriteHeader(500)
			case "id":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			case "name":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/users/get_current_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("17738362"),
			Name:  setStringPtr("sean rose"),
			Login: setStringPtr("sean@box.com"),
		},
		CreatedAt:     setTime("2012-03-26T15:43:07-07:00"),
		ModifiedAt:    setTime("2012-12-12T11:34:29-08:00"),
		Language:      setStringPtr("en"),
		SpaceAmount:   5368709120,
		SpaceUsed:     2377016,
		MaxUploadSize: 262144000,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr("Employee"),
		Phone:         setStringPtr("5555555555"),
		Address:       setStringPtr("555 Office Drive"),
		AvatarUrl:     setStringPtr("https://app.box.com/api/avatar/deprecated"),
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"normal", args{[]string{"type"}},
			normal, false, nil,
		},
		{"http error/404", args{[]string{"id"}},
			normal, true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", args{[]string{"name"}},
			normal, true, &ApiOtherError{},
		},
		{"senderror", args{[]string{"name"}},
			normal, true, &ApiOtherError{},
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

			u := NewUser(apiConn)
			got, err := u.GetCurrentUser(tt.args.fields)

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

func TestUser_GetUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		userId string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal", args{"10001", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10001",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal", args{"10002", []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10002?fields=type,id",
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
			u := NewUser(apiConn)
			got := u.GetUserReq(tt.args.userId, tt.args.fields)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}

		})
	}
}

func TestUser_GetUser(t *testing.T) {
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
			userId := strings.Split(r.URL.Path, "/")[3]

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
				w.WriteHeader(200)
				_, _ = w.Write([]byte("invalid json"))
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(200)
				resp, _ := ioutil.ReadFile("testdata/users/get_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("10543463"),
			Name:  setStringPtr("Arielle Frey"),
			Login: setStringPtr("ariellefrey@box.com"),
		},
		CreatedAt:     setTime("2011-01-07T12:37:09-08:00"),
		ModifiedAt:    setTime("2014-05-30T10:39:47-07:00"),
		Language:      setStringPtr("en"),
		Timezone:      setStringPtr("America/Los_Angeles"),
		SpaceAmount:   10737418240,
		SpaceUsed:     558732,
		MaxUploadSize: 5368709120,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr(""),
		Phone:         setStringPtr(""),
		Address:       setStringPtr(""),
		AvatarUrl:     setStringPtr("https://app.box.com/api/avatar/deprecated"),
	}
	type args struct {
		userId string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", []string{"type"}},
			normal, false, nil,
		},
		{"http error/404", args{"404", []string{"id"}},
			normal, true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", args{"999", []string{"name"}},
			normal, true, &ApiOtherError{},
		},
		{"senderror", args{"999", []string{"name"}},
			normal, true, &ApiOtherError{},
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

			u := NewUser(apiConn)
			got, err := u.GetUser(tt.args.userId, tt.args.fields)

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

func buildUserOfCommon(apiConn *APIConn) *User {
	u := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("181216415"),
			Name:  setStringPtr("sean rose"),
			Login: setStringPtr("sean+awesome@box.com"),
		},
		CreatedAt:                     setTime("2012-05-03T21:39:11-07:00"),
		ModifiedAt:                    setTime("2012-11-14T11:21:32-08:00"),
		Role:                          setUserRole(UserRoleAdmin),
		Language:                      setStringPtr("en"),
		Timezone:                      setStringPtr("Africa/Bujumbura"),
		SpaceAmount:                   11345156112,
		SpaceUsed:                     1237009912,
		MaxUploadSize:                 2147483648,
		TrackingCodes:                 []map[string]string{},
		CanSeeManagedUsers:            setBool(true),
		IsSyncEnabled:                 setBool(true),
		Status:                        setUserStatus(UserStatusActive),
		JobTitle:                      setStringPtr(""),
		Phone:                         setStringPtr("6509241374"),
		Address:                       setStringPtr(""),
		AvatarUrl:                     setStringPtr("https://www.box.com/api/avatar/large/181216415"),
		IsExemptFromDeviceLimits:      setBool(false),
		IsExemptFromLoginVerification: setBool(false),
		Enterprise: &Enterprise{
			Type: EnterpriseTypeEnterprise,
			Id:   "17077211",
			Name: "seanrose enterprise",
		},
		MyTags: &[]string{"important", "needs review"},
	}
	return u
}
func TestUser_CreateUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	u0001 := buildUserOfCommon(apiConn)
	u0001.SetLogin("10001@example.com")
	u0001.SetName("Jane Doe1")

	u0002 := buildUserOfCommon(apiConn)
	u0002.SetLogin("10002@example.com")
	u0002.SetName("Jane Doe2")
	u0002.SetRole(UserRoleCoAdmin)

	u0003 := buildUserOfCommon(apiConn)
	u0003.SetLogin("10003@example.com")
	u0003.SetName("Jane Doe3")
	u0003.SetLanguage("ja")

	u0004 := buildUserOfCommon(apiConn)
	u0004.SetLogin("10004@example.com")
	u0004.SetName("Jane Doe4")
	u0004.SetIsSyncEnabled(true)

	u0005 := buildUserOfCommon(apiConn)
	u0005.SetLogin("10005@example.com")
	u0005.SetName("Jane Doe5")
	u0005.SetJobTitle("MANAGER")

	u0006 := buildUserOfCommon(apiConn)
	u0006.SetLogin("10006@example.com")
	u0006.SetName("Jane Doe6")
	u0006.SetPhone("123-456-789")

	u0007 := buildUserOfCommon(apiConn)
	u0007.SetLogin("10007@example.com")
	u0007.SetName("Jane Doe7")
	u0007.SetAddress("1-2, ABC Street")

	u0008 := buildUserOfCommon(apiConn)
	u0008.SetLogin("10008@example.com")
	u0008.SetName("Jane Doe8")
	u0008.SetSpaceAmount(123456789)

	u0009 := buildUserOfCommon(apiConn)
	u0009.SetLogin("10009@example.com")
	u0009.SetName("Jane Doe9")
	u0009.SetTrackingCodes([]map[string]string{{"key1": "value1"}, {"key2": "value2"}})

	u0010 := buildUserOfCommon(apiConn)
	u0010.SetLogin("10010@example.com")
	u0010.SetName("Jane Doe10")
	u0010.SetCanSeeManagedUsers(false)

	u0011 := buildUserOfCommon(apiConn)
	u0011.SetLogin("10011@example.com")
	u0011.SetName("Jane Doe11")
	u0011.SetTimezone("America/Los_Angeles")

	u0012 := buildUserOfCommon(apiConn)
	u0012.SetLogin("10012@example.com")
	u0012.SetName("Jane Doe12")
	u0012.SetIsExemptFromDeviceLimits(true)

	u0013 := buildUserOfCommon(apiConn)
	u0013.SetLogin("10013@example.com")
	u0013.SetName("Jane Doe13")
	u0013.SetIsExemptFromLoginVerification(true)

	u0014 := buildUserOfCommon(apiConn)
	u0014.SetLogin("10014@example.com")
	u0014.SetName("Jane Doe14")
	u0014.SetIsExternalCollabRestricted(true)

	u0015 := buildUserOfCommon(apiConn)
	u0015.SetLogin("10015@example.com")
	u0015.SetName("Jane Doe15")
	u0015.SetStatus(UserStatusCannotDeleteEditUpload)

	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		target *User
		args   args
		want   *Request
	}{
		{"minimum",
			u0001,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10001@example.com",
	"name": "Jane Doe1"
}
`),
			},
		},
		{"role",
			u0002,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10002@example.com",
	"name": "Jane Doe2",
	"role": "coadmin"
}
`),
			},
		},
		{"language",
			u0003,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10003@example.com",
	"name": "Jane Doe3",
	"language": "ja"
}
`),
			},
		},
		{"is_sync_enabled",
			u0004,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10004@example.com",
	"name": "Jane Doe4",
	"is_sync_enabled": true
}
`),
			},
		},
		{"job_title",
			u0005,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10005@example.com",
	"name": "Jane Doe5",
	"job_title": "MANAGER"
}
`),
			},
		},
		{"phone",
			u0006,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10006@example.com",
	"name": "Jane Doe6",
	"phone": "123-456-789"
}
`),
			},
		},
		{"phone",
			u0007,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10007@example.com",
	"name": "Jane Doe7",
	"address": "1-2, ABC Street"
}
`),
			},
		},
		{"space_amount",
			u0008,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10008@example.com",
	"name": "Jane Doe8",
	"space_amount": 123456789
}
`),
			},
		},
		{"tracking_codes",
			u0009,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10009@example.com",
	"name": "Jane Doe9",
	"tracking_codes": [{"key1":"value1"},{"key2":"value2"}]
}
`),
			},
		},
		{"can_see_managed_users",
			u0010,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10010@example.com",
	"name": "Jane Doe10",
	"can_see_managed_users": false
}
`),
			},
		},
		{"timezone",
			u0011,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10011@example.com",
	"name": "Jane Doe11",
	"timezone": "America/Los_Angeles"
}
`),
			},
		},
		{"is_exempt_from_device_limits",
			u0012,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10012@example.com",
	"name": "Jane Doe12",
	"is_exempt_from_device_limits": true
}
`),
			},
		},
		{"is_exempt_from_login_verification",
			u0013,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10013@example.com",
	"name": "Jane Doe13",
	"is_exempt_from_login_verification": true
}
`),
			},
		},
		{"is_external_collab_restricted",
			u0014,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10014@example.com",
	"name": "Jane Doe14",
	"is_external_collab_restricted": true
}
`),
			},
		},
		{"status",
			u0015,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10015@example.com",
	"name": "Jane Doe15",
	"status": "cannot_delete_edit_upload"
}
`),
			},
		},
		{"fields",
			u0001,
			args{[]string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users?fields=type,id",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"login": "10001@example.com",
	"name": "Jane Doe1"
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			u := tt.target
			got := u.CreateUserReq(tt.args.fields)

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

func TestUser_CreateUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users")
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
			var v map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&v)
			if err != nil {
				t.Fatalf("there is no body data")
			}
			id := v["login"]

			switch id {
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
				resp, _ := ioutil.ReadFile("testdata/users/create_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("187273718"),
			Name:  setStringPtr("Ned Stark"),
			Login: setStringPtr("eddard@box.com"),
		},
		CreatedAt:     setTime("2012-11-15T16:34:28-08:00"),
		ModifiedAt:    setTime("2012-11-15T16:34:29-08:00"),
		Role:          setUserRole(UserRoleUser),
		Language:      setStringPtr("en"),
		Timezone:      setStringPtr("America/Los_Angeles"),
		SpaceAmount:   5368709120,
		SpaceUsed:     0,
		MaxUploadSize: 2147483648,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr(""),
		Phone:         setStringPtr("555-555-5555"),
		Address:       setStringPtr("555 Box Lane"),
		AvatarUrl:     setStringPtr("https://www.box.com/api/avatar/large/187273718"),
	}

	u1 := buildUserOfCommon(apiConn)
	u1.SetLogin("10001")
	u2 := buildUserOfCommon(apiConn)
	u2.SetLogin("404")
	u3 := buildUserOfCommon(apiConn)
	u3.SetLogin("999")
	u4 := buildUserOfCommon(apiConn)
	u4.SetLogin("999")

	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		target  *User
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"normal", u1, args{[]string{"type"}},
			normal, false, nil,
		},
		{"http error/404", u2, args{[]string{"id"}},
			nil, true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", u3, args{[]string{"name"}},
			nil, true, &ApiOtherError{},
		},
		{"senderror", u4, args{[]string{"name"}},
			nil, true, &ApiOtherError{},
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

			u := tt.target
			got, err := u.CreateUser(tt.args.fields)

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

func TestUser_UpdateUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	u0001 := buildUserOfCommon(apiConn)
	u0001.SetLogin("10001@example.com")
	u0001.SetName("Jane Doe1")

	u0002 := buildUserOfCommon(apiConn)
	u0002.SetLogin("10002@example.com")
	u0002.SetName("Jane Doe2")
	u0002.SetRole(UserRoleCoAdmin)

	u0003 := buildUserOfCommon(apiConn)
	u0003.SetLogin("10003@example.com")
	u0003.SetName("Jane Doe3")
	u0003.SetLanguage("ja")

	u0004 := buildUserOfCommon(apiConn)
	u0004.SetLogin("10004@example.com")
	u0004.SetName("Jane Doe4")
	u0004.SetIsSyncEnabled(true)

	u0005 := buildUserOfCommon(apiConn)
	u0005.SetLogin("10005@example.com")
	u0005.SetName("Jane Doe5")
	u0005.SetJobTitle("MANAGER")

	u0006 := buildUserOfCommon(apiConn)
	u0006.SetLogin("10006@example.com")
	u0006.SetName("Jane Doe6")
	u0006.SetPhone("123-456-789")

	u0007 := buildUserOfCommon(apiConn)
	u0007.SetLogin("10007@example.com")
	u0007.SetName("Jane Doe7")
	u0007.SetAddress("1-2, ABC Street")

	u0008 := buildUserOfCommon(apiConn)
	u0008.SetLogin("10008@example.com")
	u0008.SetName("Jane Doe8")
	u0008.SetSpaceAmount(123456789)

	u0009 := buildUserOfCommon(apiConn)
	u0009.SetLogin("10009@example.com")
	u0009.SetName("Jane Doe9")
	u0009.SetTrackingCodes([]map[string]string{{"key1": "value1"}, {"key2": "value2"}})

	u0010 := buildUserOfCommon(apiConn)
	u0010.SetLogin("10010@example.com")
	u0010.SetName("Jane Doe10")
	u0010.SetCanSeeManagedUsers(false)

	u0011 := buildUserOfCommon(apiConn)
	u0011.SetLogin("10011@example.com")
	u0011.SetName("Jane Doe11")
	u0011.SetTimezone("America/Los_Angeles")

	u0012 := buildUserOfCommon(apiConn)
	u0012.SetLogin("10012@example.com")
	u0012.SetName("Jane Doe12")
	u0012.SetIsExemptFromDeviceLimits(true)

	u0013 := buildUserOfCommon(apiConn)
	u0013.SetLogin("10013@example.com")
	u0013.SetName("Jane Doe13")
	u0013.SetIsExemptFromLoginVerification(true)

	u0014 := buildUserOfCommon(apiConn)
	u0014.SetLogin("10014@example.com")
	u0014.SetName("Jane Doe14")
	u0014.SetIsExternalCollabRestricted(true)

	u0015 := buildUserOfCommon(apiConn)
	u0015.SetLogin("10015@example.com")
	u0015.SetName("Jane Doe15")
	u0015.SetStatus(UserStatusCannotDeleteEditUpload)

	u0016 := buildUserOfCommon(apiConn)
	u0016.SetLogin("10016@example.com")
	u0016.SetName("Jane Doe16")
	u0016.SetIsPasswordResetRequired(true)

	u0017 := buildUserOfCommon(apiConn)
	u0017.SetLogin("10017@example.com")
	u0017.SetName("Jane Doe17")
	u0017.SetRollOutOfEnterprise(true)

	type args struct {
		userId string
		fields []string
	}
	tests := []struct {
		name   string
		target *User
		args   args
		want   *Request
	}{
		{"name",
			u0001,
			args{"10001", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10001",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe1"
}
`),
			},
		},
		{"role",
			u0002,
			args{"10002", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10002",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe2",
	"role": "coadmin"
}
`),
			},
		},
		{"language",
			u0003,
			args{"10003", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10003",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe3",
	"language": "ja"
}
`),
			},
		},
		{"is_sync_enabled",
			u0004,
			args{"10004", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10004",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe4",
	"is_sync_enabled": true
}
`),
			},
		},
		{"job_title",
			u0005,
			args{"10005", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10005",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe5",
	"job_title": "MANAGER"
}
`),
			},
		},
		{"phone",
			u0006,
			args{"10006", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10006",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe6",
	"phone": "123-456-789"
}
`),
			},
		},
		{"phone",
			u0007,
			args{"10007", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10007",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe7",
	"address": "1-2, ABC Street"
}
`),
			},
		},
		{"space_amount",
			u0008,
			args{"10008", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10008",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe8",
	"space_amount": 123456789
}
`),
			},
		},
		{"tracking_codes",
			u0009,
			args{"10009", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10009",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe9",
	"tracking_codes": [{"key1":"value1"},{"key2":"value2"}]
}
`),
			},
		},
		{"can_see_managed_users",
			u0010,
			args{"10010", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10010",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe10",
	"can_see_managed_users": false
}
`),
			},
		},
		{"timezone",
			u0011,
			args{"10011", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10011",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe11",
	"timezone": "America/Los_Angeles"
}
`),
			},
		},
		{"is_exempt_from_device_limits",
			u0012,
			args{"10012", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10012",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe12",
	"is_exempt_from_device_limits": true
}
`),
			},
		},
		{"is_exempt_from_login_verification",
			u0013,
			args{"10013", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10013",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe13",
	"is_exempt_from_login_verification": true
}
`),
			},
		},
		{"is_external_collab_restricted",
			u0014,
			args{"10014", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10014",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe14",
	"is_external_collab_restricted": true
}
`),
			},
		},
		{"status",
			u0015,
			args{"10015", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10015",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe15",
	"status": "cannot_delete_edit_upload"
}
`),
			},
		},
		{"is_password_reset_required",
			u0016,
			args{"10016", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10016",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe16",
	"is_password_reset_required": true
}
`),
			},
		},
		{"rolled out",
			u0017,
			args{"10017", nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10017",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe17",
	"enterprise": null,
	"notify": true
}
`),
			},
		},
		{"fields",
			u0001,
			args{"10001", []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10001?fields=type,id",
				Method:             PUT,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe1"
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			u := tt.target
			got := u.UpdateUserReq(tt.args.userId, tt.args.fields)

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

func TestUser_UpdateUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users")
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
			id := strings.Split(r.URL.Path, "/")[3]

			switch id {
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
				resp, _ := ioutil.ReadFile("testdata/users/create_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("187273718"),
			Name:  setStringPtr("Ned Stark"),
			Login: setStringPtr("eddard@box.com"),
		},
		CreatedAt:     setTime("2012-11-15T16:34:28-08:00"),
		ModifiedAt:    setTime("2012-11-15T16:34:29-08:00"),
		Role:          setUserRole(UserRoleUser),
		Language:      setStringPtr("en"),
		Timezone:      setStringPtr("America/Los_Angeles"),
		SpaceAmount:   5368709120,
		SpaceUsed:     0,
		MaxUploadSize: 2147483648,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr(""),
		Phone:         setStringPtr("555-555-5555"),
		Address:       setStringPtr("555 Box Lane"),
		AvatarUrl:     setStringPtr("https://www.box.com/api/avatar/large/187273718"),
	}

	u1 := buildUserOfCommon(apiConn)
	u1.SetName("10001")
	u2 := buildUserOfCommon(apiConn)
	u2.SetName("404")
	u3 := buildUserOfCommon(apiConn)
	u3.SetName("999")
	u4 := buildUserOfCommon(apiConn)
	u4.SetName("999")

	type args struct {
		userId string
		fields []string
	}
	tests := []struct {
		name    string
		target  *User
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"normal", u1, args{"10001", []string{"type"}},
			normal, false, nil,
		},
		{"http error/404", u2, args{"404", []string{"id"}},
			nil, true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", u3, args{"999", []string{"name"}},
			nil, true, &ApiOtherError{},
		},
		{"senderror", u4, args{"999", []string{"name"}},
			nil, true, &ApiOtherError{},
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

			u := tt.target
			got, err := u.UpdateUser(tt.args.userId, tt.args.fields)

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

func TestUser_CreateAppUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	u0001 := buildUserOfCommon(apiConn)
	u0001.SetLogin("10001@example.com")
	u0001.SetName("Jane Doe1")

	u0002 := buildUserOfCommon(apiConn)
	u0002.SetLogin("10002@example.com")
	u0002.SetName("Jane Doe2")
	u0002.SetRole(UserRoleCoAdmin)

	u0003 := buildUserOfCommon(apiConn)
	u0003.SetLogin("10003@example.com")
	u0003.SetName("Jane Doe3")
	u0003.SetLanguage("ja")

	u0004 := buildUserOfCommon(apiConn)
	u0004.SetLogin("10004@example.com")
	u0004.SetName("Jane Doe4")
	u0004.SetIsSyncEnabled(true)

	u0005 := buildUserOfCommon(apiConn)
	u0005.SetLogin("10005@example.com")
	u0005.SetName("Jane Doe5")
	u0005.SetJobTitle("MANAGER")

	u0006 := buildUserOfCommon(apiConn)
	u0006.SetLogin("10006@example.com")
	u0006.SetName("Jane Doe6")
	u0006.SetPhone("123-456-789")

	u0007 := buildUserOfCommon(apiConn)
	u0007.SetLogin("10007@example.com")
	u0007.SetName("Jane Doe7")
	u0007.SetAddress("1-2, ABC Street")

	u0008 := buildUserOfCommon(apiConn)
	u0008.SetLogin("10008@example.com")
	u0008.SetName("Jane Doe8")
	u0008.SetSpaceAmount(123456789)

	u0009 := buildUserOfCommon(apiConn)
	u0009.SetLogin("10009@example.com")
	u0009.SetName("Jane Doe9")
	u0009.SetTrackingCodes([]map[string]string{{"key1": "value1"}, {"key2": "value2"}})

	u0010 := buildUserOfCommon(apiConn)
	u0010.SetLogin("10010@example.com")
	u0010.SetName("Jane Doe10")
	u0010.SetCanSeeManagedUsers(false)

	u0011 := buildUserOfCommon(apiConn)
	u0011.SetLogin("10011@example.com")
	u0011.SetName("Jane Doe11")
	u0011.SetTimezone("America/Los_Angeles")

	u0012 := buildUserOfCommon(apiConn)
	u0012.SetLogin("10012@example.com")
	u0012.SetName("Jane Doe12")
	u0012.SetIsExemptFromDeviceLimits(true)

	u0013 := buildUserOfCommon(apiConn)
	u0013.SetLogin("10013@example.com")
	u0013.SetName("Jane Doe13")
	u0013.SetIsExemptFromLoginVerification(true)

	u0014 := buildUserOfCommon(apiConn)
	u0014.SetLogin("10014@example.com")
	u0014.SetName("Jane Doe14")
	u0014.SetIsExternalCollabRestricted(true)

	u0015 := buildUserOfCommon(apiConn)
	u0015.SetLogin("10015@example.com")
	u0015.SetName("Jane Doe15")
	u0015.SetStatus(UserStatusCannotDeleteEditUpload)

	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		target *User
		args   args
		want   *Request
	}{
		{"minimum",
			u0001,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe1",
	"is_platform_access_only": true
}
`),
			},
		},
		{"role",
			u0002,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe2",
	"is_platform_access_only": true
}
`),
			},
		},
		{"language",
			u0003,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe3",
	"language": "ja",
	"is_platform_access_only": true
}
`),
			},
		},
		{"is_sync_enabled",
			u0004,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe4",
	"is_platform_access_only": true
}
`),
			},
		},
		{"job_title",
			u0005,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe5",
	"job_title": "MANAGER",
	"is_platform_access_only": true
}
`),
			},
		},
		{"phone",
			u0006,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe6",
	"phone": "123-456-789",
	"is_platform_access_only": true
}
`),
			},
		},
		{"address",
			u0007,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe7",
	"address": "1-2, ABC Street",
	"is_platform_access_only": true
}
`),
			},
		},
		{"space_amount",
			u0008,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe8",
	"space_amount": 123456789,
	"is_platform_access_only": true
}
`),
			},
		},
		{"tracking_codes",
			u0009,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe9",
	"is_platform_access_only": true
}
`),
			},
		},
		{"can_see_managed_users",
			u0010,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe10",
	"can_see_managed_users": false,
	"is_platform_access_only": true
}
`),
			},
		},
		{"timezone",
			u0011,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe11",
	"timezone": "America/Los_Angeles",
	"is_platform_access_only": true
}
`),
			},
		},
		{"is_exempt_from_device_limits",
			u0012,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe12",
	"is_platform_access_only": true
}
`),
			},
		},
		{"is_exempt_from_login_verification",
			u0013,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe13",
	"is_platform_access_only": true
}
`),
			},
		},
		{"is_external_collab_restricted",
			u0014,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe14",
	"is_external_collab_restricted": true,
	"is_platform_access_only": true
}
`),
			},
		},
		{"status",
			u0015,
			args{nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe15",
	"status": "cannot_delete_edit_upload",
	"is_platform_access_only": true
}
`),
			},
		},
		{"fields",
			u0001,
			args{[]string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users?fields=type,id",
				Method:             POST,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body: strings.NewReader(`
{
	"name": "Jane Doe1",
	"is_platform_access_only": true
}
`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			u := tt.target
			got := u.CreateAppUserReq(tt.args.fields)

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

func TestUser_CreateAppUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users")
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
			var v map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&v)
			if err != nil {
				t.Fatalf("there is no body data")
			}
			id := v["name"]

			switch id {
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
				resp, _ := ioutil.ReadFile("testdata/users/create_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	normal := &User{
		apiInfo: &apiInfo{api: apiConn},
		UserGroupMini: UserGroupMini{
			Type:  setUserType(TYPE_USER),
			ID:    setStringPtr("187273718"),
			Name:  setStringPtr("Ned Stark"),
			Login: setStringPtr("eddard@box.com"),
		},
		CreatedAt:     setTime("2012-11-15T16:34:28-08:00"),
		ModifiedAt:    setTime("2012-11-15T16:34:29-08:00"),
		Role:          setUserRole(UserRoleUser),
		Language:      setStringPtr("en"),
		Timezone:      setStringPtr("America/Los_Angeles"),
		SpaceAmount:   5368709120,
		SpaceUsed:     0,
		MaxUploadSize: 2147483648,
		Status:        setUserStatus(UserStatusActive),
		JobTitle:      setStringPtr(""),
		Phone:         setStringPtr("555-555-5555"),
		Address:       setStringPtr("555 Box Lane"),
		AvatarUrl:     setStringPtr("https://www.box.com/api/avatar/large/187273718"),
	}

	u1 := buildUserOfCommon(apiConn)
	u1.SetName("10001")
	u2 := buildUserOfCommon(apiConn)
	u2.SetName("404")
	u3 := buildUserOfCommon(apiConn)
	u3.SetName("999")
	u4 := buildUserOfCommon(apiConn)
	u4.SetName("999")

	type args struct {
		fields []string
	}
	tests := []struct {
		name    string
		target  *User
		args    args
		want    *User
		wantErr bool
		errType interface{}
	}{
		{"normal", u1, args{[]string{"type"}},
			normal, false, nil,
		},
		{"http error/404", u2, args{[]string{"id"}},
			nil, true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", u3, args{[]string{"name"}},
			nil, true, &ApiOtherError{},
		},
		{"senderror", u4, args{[]string{"name"}},
			nil, true, &ApiOtherError{},
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

			u := tt.target
			got, err := u.CreateAppUser(tt.args.fields)

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

func TestUser_DeleteUserReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		userId string
		notify bool
		force  bool
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"normal",
			args{"10001", true, true},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10001?notify=true&force=true",
				Method:             DELETE,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10002", false, true},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10002?notify=false&force=true",
				Method:             DELETE,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10003", false, false},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10003?notify=false&force=false",
				Method:             DELETE,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"normal",
			args{"10004", true, false},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users/10004?notify=true&force=false",
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

			u := NewUser(apiConn)
			got := u.DeleteUserReq(tt.args.userId, tt.args.notify, tt.args.force)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestUser_DeleteUser(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users")
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
			id := strings.Split(r.URL.Path, "/")[3]

			switch id {
			case "500":
				w.WriteHeader(500)
			case "404":
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(404)
				resp, _ := ioutil.ReadFile("testdata/genericerror/404.json")
				_, _ = w.Write(resp)
			default:
				w.Header().Set("content-Type", "application/json")
				w.WriteHeader(204)
				resp, _ := ioutil.ReadFile("testdata/users/create_user.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	u1 := buildUserOfCommon(apiConn)
	u1.SetName("10001")
	u2 := buildUserOfCommon(apiConn)
	u2.SetName("404")
	u3 := buildUserOfCommon(apiConn)
	u3.SetName("999")

	type args struct {
		userId string
		notify bool
		force  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errType interface{}
	}{
		{"normal", args{"10001", false, false},
			false, nil,
		},
		{"http error/404", args{"404", false, false},
			true, &ApiStatusError{Status: 404},
		},
		{"senderror", args{"999", false, false},
			true, &ApiOtherError{},
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

			u := NewUser(apiConn)
			err := u.DeleteUser(tt.args.userId, tt.args.notify, tt.args.force)

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

func TestUser_GetEnterpriseUsersReq(t *testing.T) {
	url := "https://example.com"
	apiConn := commonInit(url)

	type args struct {
		filterTerm string
		offset     int
		limit      int
		fields     []string
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"filterTerm empty",
			args{"", 0, 1001, nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users?offset=0&limit=1000",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"filterTerm",
			args{"あいうえお", 1, 999, nil},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users?offset=1&limit=999&filter_term=%E3%81%82%E3%81%84%E3%81%86%E3%81%88%E3%81%8A",
				Method:             GET,
				shouldAuthenticate: true,
				numRedirects:       defaultNumRedirects,
				headers:            http.Header{},
				body:               nil,
			},
		},
		{"filterTerm",
			args{"something", 1, 999, []string{"type", "id"}},
			&Request{
				apiConn:            apiConn,
				Url:                url + "/2.0/users?offset=1&limit=999&filter_term=something&fields=type,id",
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

			u := NewUser(apiConn)
			got := u.GetEnterpriseUsersReq(tt.args.filterTerm, tt.args.offset, tt.args.limit, tt.args.fields)

			opts := diffCompOptions(*got)
			opts = append(opts, cmpopts.IgnoreUnexported(Request{}))

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("diff:  (-got +want)\n%s", diff)
				return
			}
		})
	}
}

func TestUser_GetEnterpriseUsers(t *testing.T) {
	// test server (dummy box api)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// for request.Send() return error (auth failed)
			if strings.HasPrefix(r.URL.Path, "/oauth2/token") {
				w.WriteHeader(401)
				return
			}
			// URL check
			if !strings.HasPrefix(r.URL.Path, "/2.0/users") {
				t.Errorf("invalid access url %s : %s", r.URL.Path, "/2.0/users")
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
			id := r.URL.Query().Get("filter_term")

			switch id {
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
				resp, _ := ioutil.ReadFile("testdata/users/get_enterprise_users.json")
				_, _ = w.Write(resp)
			}
			return
		},
	))
	defer ts.Close()

	apiConn := commonInit(ts.URL)

	type args struct {
		filterTerm string
		offset     int
		limit      int
		fields     []string
	}
	tests := []struct {
		name           string
		args           args
		want           []*User
		wantOffset     int
		wantLimit      int
		wantTotalCount int
		wantErr        bool
		errType        interface{}
	}{
		{"normal", args{"10001", 0, 1000, nil},
			[]*User{
				{
					apiInfo: &apiInfo{api: apiConn},
					UserGroupMini: UserGroupMini{
						Type:  setUserType(TYPE_USER),
						ID:    setStringPtr("181216415"),
						Name:  setStringPtr("sean rose"),
						Login: setStringPtr("sean+awesome@box.com"),
					},
					CreatedAt:     setTime("2012-05-03T21:39:11-07:00"),
					ModifiedAt:    setTime("2012-08-23T14:57:48-07:00"),
					Language:      setStringPtr("en"),
					SpaceAmount:   5368709120,
					SpaceUsed:     52947,
					MaxUploadSize: 104857600,
					Status:        setUserStatus(UserStatusActive),
					JobTitle:      setStringPtr(""),
					Phone:         setStringPtr("5555551374"),
					Address:       setStringPtr("10 Cloud Way Los Altos CA"),
					AvatarUrl:     setStringPtr("https://app.box.com/api/avatar/large/181216415"),
				},
			},
			0, 1000, 1, false, nil,
		},
		{"http error/404", args{"404", 0, 1000, nil},
			[]*User{},
			0, 0, 0,
			true, &ApiStatusError{Status: 404},
		},
		{"returned invalid json/999", args{"999", 0, 1000, []string{"name"}},
			nil, 0, 0, 0,
			true, &ApiOtherError{},
		},
		{"senderror", args{"999", 0, 1000, nil},
			nil, 0, 0, 0,
			true, &ApiOtherError{},
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

			u := NewUser(apiConn)
			got, gotOffset, gotLimit, gotTotalCount, err := u.GetEnterpriseUsers(tt.args.filterTerm, tt.args.offset, tt.args.limit, tt.args.fields)

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

			if gotOffset != 0 {

			}
			if gotLimit != tt.wantLimit || gotOffset != tt.wantOffset || gotTotalCount != tt.wantTotalCount {
				t.Errorf("returned offset/limit/totalCount was incorrect")
				return
			}

			// If normal response
			opts := diffCompOptions(User{}, apiInfo{})
			if diff := cmp.Diff(&got, &tt.want, opts...); diff != "" {
				t.Errorf("Marshal/Unmarshal differs: (-got +want)\n%s", diff)
				return
			}
			// exists apiInfo
			for _, v := range got {
				if v.apiInfo == nil {
					t.Errorf("not exists `apiInfo` field\n")
					return
				}
			}

		})
	}
}
