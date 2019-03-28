package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
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
		// TODO: Add test cases.
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

func TestUserRole_String(t *testing.T) {
	tests := []struct {
		name string
		ur   *UserRole
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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

func TestEnterpriseType_String(t *testing.T) {
	tests := []struct {
		name string
		et   *EnterpriseType
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
