package goboxer

import (
	"reflect"
	"testing"
)

func TestUserGroupMini_IsUser(t *testing.T) {
	type fields struct {
		Type  *UserGroupType
		ID    *string
		Name  *string
		Login *string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"nil", fields{nil, nil, nil, nil}, false},
		{"user",
			fields{
				setUserType(TYPE_USER),
				setStringPtr("1"),
				setStringPtr("Name1"),
				setStringPtr("login1@example.com"),
			},
			true,
		},
		{"group",
			fields{
				setUserType(TYPE_GROUP),
				setStringPtr("2"),
				setStringPtr("Name2"),
				setStringPtr("login2@example.com"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *UserGroupMini
			if tt.fields.Type == nil {
				u = nil
			} else {
				u = &UserGroupMini{
					Type:  tt.fields.Type,
					ID:    tt.fields.ID,
					Name:  tt.fields.Name,
					Login: tt.fields.Login,
				}
			}
			if got := u.IsUser(); got != tt.want {
				t.Errorf("UserGroupMini.IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGroupMini_String(t *testing.T) {
	type fields struct {
		Type  *UserGroupType
		ID    *string
		Name  *string
		Login *string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"nil", fields{nil, nil, nil, nil}, "<nil>"},
		{"normal",
			fields{
				setUserType(TYPE_USER),
				setStringPtr("1"),
				setStringPtr("Name1"),
				setStringPtr("login1@example.com"),
			},
			`{ Type:user, ID:1, Name:Name1, Login:login1@example.com }`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u *UserGroupMini
			if tt.fields.Type == nil {
				u = nil
			} else {
				u = &UserGroupMini{
					Type:  tt.fields.Type,
					ID:    tt.fields.ID,
					Name:  tt.fields.Name,
					Login: tt.fields.Login,
				}
			}
			if got := u.String(); got != tt.want {
				t.Errorf("UserGroupMini.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGroupType_String(t *testing.T) {
	tests := []struct {
		name string
		u    *UserGroupType
		want string
	}{
		{"nil", nil, "<nil>"},
		{"user", setUserType(TYPE_USER), "user"},
		{"group", setUserType(TYPE_GROUP), "group"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("UserGroupType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserGroupType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       *UserGroupType
		want    []byte
		wantErr bool
	}{
		{"nil", nil, []byte(`null`), false},
		{"user", setUserType(TYPE_USER), []byte(`"user"`), false},
		{"group", setUserType(TYPE_GROUP), []byte(`"group"`), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserGroupType.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserGroupType.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
