package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type InvitabilityLevel string

const (
	InvitabilityAdminsOnly      InvitabilityLevel = "admins_only"
	InvitabilityAdminsMembers   InvitabilityLevel = "admins_and_members"
	InvitabilityAllManagedUsers InvitabilityLevel = "all_managed_users"
)

func (us *InvitabilityLevel) String() string {
	if us == nil {
		return "<nil>"
	}
	return string(*us)
}

func (us *InvitabilityLevel) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + us.String() + `"`), nil
	}
}

type MemberViewabilityLevel string

const (
	MemberViewabilityAdminsOnly      MemberViewabilityLevel = "admins_only"
	MemberViewabilityAdminsMembers   MemberViewabilityLevel = "admins_and_members"
	MemberViewabilityAllManagedUsers MemberViewabilityLevel = "all_managed_users"
)

func (us *MemberViewabilityLevel) String() string {
	if us == nil {
		return "<nil>"
	}
	return string(*us)
}

func (us *MemberViewabilityLevel) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + us.String() + `"`), nil
	}
}

type Group struct {
	UserGroupMini
	apiInfo                *apiInfo                `json:"-"`
	CreatedAt              *time.Time              `json:"created_at,omitempty"`
	ModifiedAt             *time.Time              `json:"modified_at,omitempty"`
	Provenance             *string                 `json:"provenance,omitempty"`
	ExternalSyncIdentifier *string                 `json:"external_sync_identifier,omitempty"`
	Description            *string                 `json:"description,omitempty"`
	InvitabilityLevel      *InvitabilityLevel      `json:"invitability_level,omitempty"`
	MemberViewabilityLevel *MemberViewabilityLevel `json:"member_viewability_level,omitempty"`
}

func (g *Group) Type() string {
	return g.UserGroupMini.Type.String()
}

func NewGroup(api *ApiConn) *Group {
	return &Group{
		apiInfo: &apiInfo{api: api},
	}
}

var GroupAllFields = []string{
	"type", "id", "name", "created_at", "modified_at",
	"provenance", "external_sync_identifier", "description",
	"invitability_level", "member_viewability_level",
}

func (g *Group) GetGroupReq(groupId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", g.apiInfo.api.BaseURL, "groups/", groupId, BuildFieldsQueryParams(fields))

	return NewRequest(g.apiInfo.api, url, GET, nil, nil)
}
func (g *Group) GetGroup(groupId string, fields []string) (*Group, error) {

	req := g.GetGroupReq(groupId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get group info"))
		return nil, err
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Group) CreateGroupReq(fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s?%s", g.apiInfo.api.BaseURL, "groups", BuildFieldsQueryParams(fields))

	b, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
	}
	return NewRequest(g.apiInfo.api, url, POST, nil, bytes.NewReader(b))
}
func (g *Group) CreateGroup(fields []string) (*Group, error) {

	req := g.CreateGroupReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to create group"))
		return nil, err
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Group) SetName(name string) *Group {
	g.Name = &name
	return g
}
func (g *Group) SetProvenance(provenance string) *Group {
	g.Provenance = &provenance
	return g
}
func (g *Group) SetExternalSyncIdentifier(externalSyncIdentifier string) *Group {
	g.ExternalSyncIdentifier = &externalSyncIdentifier
	return g
}
func (g *Group) SetDescription(description string) *Group {
	g.Description = &description
	return g
}
func (g *Group) SetInvitabilityLevel(invitabilityLevel InvitabilityLevel) *Group {
	g.InvitabilityLevel = &invitabilityLevel
	return g
}
func (g *Group) SetMemberViewabiityLevel(memberViewabilityLevel MemberViewabilityLevel) *Group {
	g.MemberViewabilityLevel = &memberViewabilityLevel
	return g
}

func (g *Group) UpdateGroupReq(groupId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", g.apiInfo.api.BaseURL, "groups/", groupId, BuildFieldsQueryParams(fields))

	b, err := json.Marshal(g)
	if err != nil {
		fmt.Println(err)
	}
	return NewRequest(g.apiInfo.api, url, PUT, nil, bytes.NewReader(b))
}
func (g *Group) UpdateGroup(groupId string, fields []string) (*Group, error) {

	req := g.UpdateGroupReq(groupId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to update group info"))
		return nil, err
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Group) DeleteGroupReq(groupId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", g.apiInfo.api.BaseURL, "groups/", groupId)

	return NewRequest(g.apiInfo.api, url, DELETE, nil, nil)
}
func (g *Group) DeleteGroup(groupId string) error {

	req := g.DeleteGroupReq(groupId)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to delete group"))
		return err
	}
	return nil
}

func (g *Group) GetEnterpriseGroupsReq(name string, offset int32, limit int32, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s?&name=%s&offset=%d&limit=%d&%s", g.apiInfo.api.BaseURL, "groups", name, offset, limit, BuildFieldsQueryParams(fields))

	return NewRequest(g.apiInfo.api, url, GET, nil, nil)
}
func (g *Group) GetEnterpriseGroups(name string, offset int32, limit int32, fields []string) (outGroups []*Group, outOffset int, outLimit int, outTotalCount int, err error) {

	req := g.GetEnterpriseGroupsReq(name, offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get enterprise groups"))
		return nil, 0, 0, 0, err
	}

	groups := struct {
		TotalCount int      `json:"total_count"`
		Entries    []*Group `json:"entries"`
		Offset     int      `json:"offset"`
		Limit      int      `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &groups)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return groups.Entries, groups.Offset, groups.Limit, groups.TotalCount, nil
}
