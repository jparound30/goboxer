package goboxer

import "fmt"

type UserGroupMini struct {
	Type  *UserGroupType `json:"type,omitempty"`
	ID    *string        `json:"id,omitempty"`
	Name  *string        `json:"name,omitempty"`
	Login *string        `json:"login,omitempty"`
}

func (u *UserGroupMini) IsUser() bool {
	if u == nil {
		return false
	} else if *u.Type == TYPE_USER {
		return true
	} else {
		return false
	}
}

func (u *UserGroupMini) String() string {
	if u == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ Type:%s, ID:%s, Name:%s, Login:%s }", u.Type.String(), toString(u.ID), toString(u.Name), toString(u.Login))
}

type UserGroupType string

func (u *UserGroupType) String() string {
	if u == nil {
		return "<nil>"
	}
	return string(*u)
}
func (u *UserGroupType) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + u.String() + `"`), nil
	}
}

const (
	TYPE_USER  UserGroupType = "user"
	TYPE_GROUP UserGroupType = "group"
)
