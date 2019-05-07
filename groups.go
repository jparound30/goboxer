package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	apiInfo                *apiInfo
	CreatedAt              *time.Time              `json:"created_at,omitempty"`
	ModifiedAt             *time.Time              `json:"modified_at,omitempty"`
	Provenance             *string                 `json:"provenance,omitempty"`
	ExternalSyncIdentifier *string                 `json:"external_sync_identifier,omitempty"`
	Description            *string                 `json:"description,omitempty"`
	InvitabilityLevel      *InvitabilityLevel      `json:"invitability_level,omitempty"`
	MemberViewabilityLevel *MemberViewabilityLevel `json:"member_viewability_level,omitempty"`

	changeFlag uint64
}

const (
	cGroupName uint64 = 1 << (iota)
	cGroupProvenance
	cGroupExternalSyncIdentifier
	cGroupDescription
	cGroupInvitabilityLevel
	cGroupMemberViewabilityLevel
)

func (g *Group) ResourceType() BoxResourceType {
	return GroupResource
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

// Get Group
//
// Get information about a group.
// https://developer.box.com/reference#get-group
func (g *Group) GetGroupReq(groupId string, fields []string) *Request {
	var baseUrl string
	var query string
	baseUrl = fmt.Sprintf("%s%s%s", g.apiInfo.api.BaseURL, "groups/", groupId)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}
	return NewRequest(g.apiInfo.api, baseUrl+query, GET, nil, nil)
}

// Get Group
//
// Get information about a group.
// https://developer.box.com/reference#get-group
func (g *Group) GetGroup(groupId string, fields []string) (*Group, error) {

	req := g.GetGroupReq(groupId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Create Group
//
// Create a new group. Only admin roles can create and manage groups.
// https://developer.box.com/reference#create-a-group
func (g *Group) CreateGroupReq(fields []string) *Request {
	var baseUrl string
	var query string

	baseUrl = fmt.Sprintf("%s%s", g.apiInfo.api.BaseURL, "groups")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &Group{}
	if g.changeFlag&cGroupName == cGroupName {
		data.Name = g.Name
	}
	if g.changeFlag&cGroupProvenance == cGroupProvenance {
		data.Provenance = g.Provenance
	}
	if g.changeFlag&cGroupExternalSyncIdentifier == cGroupExternalSyncIdentifier {
		data.ExternalSyncIdentifier = g.ExternalSyncIdentifier
	}
	if g.changeFlag&cGroupDescription == cGroupDescription {
		data.Description = g.Description
	}
	if g.changeFlag&cGroupInvitabilityLevel == cGroupInvitabilityLevel {
		data.InvitabilityLevel = g.InvitabilityLevel
	}
	if g.changeFlag&cGroupMemberViewabilityLevel == cGroupMemberViewabilityLevel {
		data.MemberViewabilityLevel = g.MemberViewabilityLevel
	}

	b, _ := json.Marshal(data)
	return NewRequest(g.apiInfo.api, baseUrl+query, POST, nil, bytes.NewReader(b))
}

// Create Group
//
// Create a new group. Only admin roles can create and manage groups.
// https://developer.box.com/reference#create-a-group
func (g *Group) CreateGroup(fields []string) (*Group, error) {

	req := g.CreateGroupReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Group) SetName(name string) *Group {
	g.Name = &name
	g.changeFlag |= cGroupName
	return g
}
func (g *Group) SetProvenance(provenance string) *Group {
	g.Provenance = &provenance
	g.changeFlag |= cGroupProvenance
	return g
}
func (g *Group) SetExternalSyncIdentifier(externalSyncIdentifier string) *Group {
	g.ExternalSyncIdentifier = &externalSyncIdentifier
	g.changeFlag |= cGroupExternalSyncIdentifier
	return g
}
func (g *Group) SetDescription(description string) *Group {
	g.Description = &description
	g.changeFlag |= cGroupDescription
	return g
}
func (g *Group) SetInvitabilityLevel(invitabilityLevel InvitabilityLevel) *Group {
	g.InvitabilityLevel = &invitabilityLevel
	g.changeFlag |= cGroupInvitabilityLevel
	return g
}
func (g *Group) SetMemberViewabiityLevel(memberViewabilityLevel MemberViewabilityLevel) *Group {
	g.MemberViewabilityLevel = &memberViewabilityLevel
	g.changeFlag |= cGroupMemberViewabilityLevel
	return g
}

// Update Group
//
// Update a group.
// https://developer.box.com/reference#update-a-group
func (g *Group) UpdateGroupReq(groupId string, fields []string) *Request {
	var baseUrl string
	var query string
	baseUrl = fmt.Sprintf("%s%s%s", g.apiInfo.api.BaseURL, "groups/", groupId)

	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &Group{}
	if g.changeFlag&cGroupName == cGroupName {
		data.Name = g.Name
	}
	if g.changeFlag&cGroupProvenance == cGroupProvenance {
		data.Provenance = g.Provenance
	}
	if g.changeFlag&cGroupExternalSyncIdentifier == cGroupExternalSyncIdentifier {
		data.ExternalSyncIdentifier = g.ExternalSyncIdentifier
	}
	if g.changeFlag&cGroupDescription == cGroupDescription {
		data.Description = g.Description
	}
	if g.changeFlag&cGroupInvitabilityLevel == cGroupInvitabilityLevel {
		data.InvitabilityLevel = g.InvitabilityLevel
	}
	if g.changeFlag&cGroupMemberViewabilityLevel == cGroupMemberViewabilityLevel {
		data.MemberViewabilityLevel = g.MemberViewabilityLevel
	}

	b, _ := json.Marshal(data)
	return NewRequest(g.apiInfo.api, baseUrl+query, PUT, nil, bytes.NewReader(b))
}

// Update Group
//
// Update a group.
// https://developer.box.com/reference#update-a-group
func (g *Group) UpdateGroup(groupId string, fields []string) (*Group, error) {

	req := g.UpdateGroupReq(groupId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &Group{apiInfo: &apiInfo{api: g.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete Group
//
// Delete a group.
// //https://developer.box.com/reference#delete-a-group
func (g *Group) DeleteGroupReq(groupId string) *Request {
	var baseUrl string
	baseUrl = fmt.Sprintf("%s%s%s", g.apiInfo.api.BaseURL, "groups/", groupId)

	return NewRequest(g.apiInfo.api, baseUrl, DELETE, nil, nil)
}

// Delete Group
//
// Delete a group.
// //https://developer.box.com/reference#delete-a-group
func (g *Group) DeleteGroup(groupId string) error {

	req := g.DeleteGroupReq(groupId)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		return newApiStatusError(resp.Body)
	}
	return nil
}

// Get Enterprise Groups
//
// Returns all of the groups for given enterprise. Must have permissions to see an enterprise's groups.
// https://developer.box.com/reference#groups
func (g *Group) GetEnterpriseGroupsReq(name string, offset int32, limit int32, fields []string) *Request {
	var urlBase string
	var query string
	urlBase = fmt.Sprintf("%s%s", g.apiInfo.api.BaseURL, "groups")
	if limit > 1000 {
		limit = 1000
	}
	query = fmt.Sprintf("?offset=%d&limit=%d", offset, limit)
	if name != "" {
		query += fmt.Sprintf("&name=%s", url.QueryEscape(name))
	}
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query += fmt.Sprintf("&%s", fieldsParams)
	}
	return NewRequest(g.apiInfo.api, urlBase+query, GET, nil, nil)
}

// Get Enterprise Groups
//
// Returns all of the groups for given enterprise. Must have permissions to see an enterprise's groups.
// https://developer.box.com/reference#groups
func (g *Group) GetEnterpriseGroups(name string, offset int32, limit int32, fields []string) (outGroups []*Group, outOffset int, outLimit int, outTotalCount int, err error) {

	req := g.GetEnterpriseGroupsReq(name, offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, 0, 0, 0, newApiStatusError(resp.Body)
	}

	groups := struct {
		TotalCount int      `json:"total_count"`
		Entries    []*Group `json:"entries"`
		Offset     int      `json:"offset"`
		Limit      int      `json:"limit"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &groups)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return groups.Entries, groups.Offset, groups.Limit, groups.TotalCount, nil
}
