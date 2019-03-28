package goboxer

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

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
		// TODO: Add test cases.
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
