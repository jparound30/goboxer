package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MembershipRole string

const (
	MembershipRoleAdmin  MembershipRole = "admin"
	MembershipRoleMember MembershipRole = "member"
)

func (us *MembershipRole) String() string {
	if us == nil {
		return "<nil>"
	}
	return string(*us)
}

func (us *MembershipRole) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + us.String() + `"`), nil
	}
}

type ConfigurablePermissions struct {
	CanRunReports     *bool `json:"can_run_reports,omitempty"`
	CanInstantLogin   *bool `json:"can_instant_login,omitempty"`
	CanCreateAccounts *bool `json:"can_create_accounts,omitempty"`
	CanEditAccounts   *bool `json:"can_edit_accounts,omitempty"`
}

type Membership struct {
	apiInfo                 *apiInfo
	Type                    *string                  `json:"type,omitempty"`
	ID                      *string                  `json:"id,omitempty"`
	User                    *UserGroupMini           `json:"user,omitempty"`
	Group                   *UserGroupMini           `json:"group,omitempty"`
	Role                    *MembershipRole          `json:"role,omitempty"`
	CreatedAt               *time.Time               `json:"created_at,omitempty"`
	ModifiedAt              *time.Time               `json:"modified_at,omitempty"`
	ConfigurablePermissions *ConfigurablePermissions `json:"configurable_permissions,omitempty"`
}

func (m *Membership) ResourceType() BoxResourceType {
	return MembershipResource
}

func NewMembership(api *ApiConn) *Membership {
	return &Membership{
		apiInfo: &apiInfo{api: api},
	}
}

// Get Membership
//
// Fetches a specific group membership entry.
// https://developer.box.com/reference#get-a-group-membership-entry
func (m *Membership) GetMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships/", membershipId)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}

// Get Membership
//
// Fetches a specific group membership entry.
// https://developer.box.com/reference#get-a-group-membership-entry
func (m *Membership) GetMembership(membershipId string) (*Membership, error) {

	req := m.GetMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	membership := Membership{apiInfo: m.apiInfo}

	err = UnmarshalJSONWrapper(resp.Body, &membership)
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// Create Membership
//
// Add a member to a group.
// https://developer.box.com/reference#add-a-member-to-a-group
func (m *Membership) CreateMembershipReq() *Request {
	var url string
	url = fmt.Sprintf("%s%s", m.apiInfo.api.BaseURL, "group_memberships")
	data := &Membership{
		User:  &UserGroupMini{ID: m.User.ID},
		Group: &UserGroupMini{ID: m.Group.ID},
	}
	if m.ConfigurablePermissions != nil {
		data.ConfigurablePermissions = m.ConfigurablePermissions
	}
	if m.Role != nil {
		data.Role = m.Role
	}
	b, _ := json.Marshal(data)
	return NewRequest(m.apiInfo.api, url, POST, nil, bytes.NewReader(b))
}

// Create Membership
//
// Add a member to a group.
// https://developer.box.com/reference#add-a-member-to-a-group
func (m *Membership) CreateMembership() (*Membership, error) {

	req := m.CreateMembershipReq()
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	membership := &Membership{apiInfo: m.apiInfo}

	err = UnmarshalJSONWrapper(resp.Body, membership)
	if err != nil {
		return nil, err
	}
	return membership, nil
}

func (m *Membership) SetUser(userId string) *Membership {
	m.User = &UserGroupMini{
		ID: &userId,
	}
	return m
}
func (m *Membership) SetGroup(groupId string) *Membership {
	m.Group = &UserGroupMini{
		ID: &groupId,
	}
	return m
}
func (m *Membership) SetRole(role MembershipRole) *Membership {
	m.Role = &role
	return m
}
func (m *Membership) SetConfigurablePermissions(canRunReports, canInstantLogin, canCreateAccounts, canEditAccounts bool) *Membership {
	m.ConfigurablePermissions = &ConfigurablePermissions{
		CanRunReports:     &canRunReports,
		CanInstantLogin:   &canInstantLogin,
		CanCreateAccounts: &canCreateAccounts,
		CanEditAccounts:   &canEditAccounts,
	}
	return m
}

// Update Membership
//
// Update a group membership.
// https://developer.box.com/reference#update-a-group-membership
func (m *Membership) UpdateMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships/", membershipId)

	data := &Membership{}
	if m.ConfigurablePermissions != nil {
		data.ConfigurablePermissions = m.ConfigurablePermissions
	}
	if m.Role != nil {
		data.Role = m.Role
	}
	b, _ := json.Marshal(data)
	return NewRequest(m.apiInfo.api, url, PUT, nil, bytes.NewReader(b))
}

// Update Membership
//
// Update a group membership.
// https://developer.box.com/reference#update-a-group-membership
func (m *Membership) UpdateMembership(membershipId string) (*Membership, error) {

	req := m.UpdateMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	membership := Membership{apiInfo: m.apiInfo}

	err = UnmarshalJSONWrapper(resp.Body, &membership)
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// Delete Membership
//
// Delete a group membership.
// https://developer.box.com/reference#delete-a-group-membership
func (m *Membership) DeleteMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships/", membershipId)

	return NewRequest(m.apiInfo.api, url, DELETE, nil, nil)
}

// Delete Membership
//
// Delete a group membership.
// https://developer.box.com/reference#delete-a-group-membership
func (m *Membership) DeleteMembership(membershipId string) error {

	req := m.DeleteMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		return newApiStatusError(resp.Body)
	}
	return nil
}

// Get Memberships for Group
//
// Returns all of the members for a given group if the requesting user has access.
// https://developer.box.com/reference#get-the-membership-list-for-a-group
func (m *Membership) GetMembershipForGroupReq(groupId string, offset int32, limit int32) *Request {
	var url string
	if limit > 1000 {
		limit = 1000
	}
	url = fmt.Sprintf("%s%s%s%s?offset=%d&limit=%d", m.apiInfo.api.BaseURL, "groups/", groupId, "/memberships", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}

// Get Memberships for Group
//
// Returns all of the members for a given group if the requesting user has access.
// https://developer.box.com/reference#get-the-membership-list-for-a-group
func (m *Membership) GetMembershipForGroup(groupId string, offset int32, limit int32) (outMembership []*Membership, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetMembershipForGroupReq(groupId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, 0, 0, 0, newApiStatusError(resp.Body)
	}

	memberships := struct {
		TotalCount int           `json:"total_count"`
		Entries    []*Membership `json:"entries"`
		Offset     int           `json:"offset"`
		Limit      int           `json:"limit"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &memberships)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	for _, v := range memberships.Entries {
		v.apiInfo = m.apiInfo
	}
	return memberships.Entries, memberships.Offset, memberships.Limit, memberships.TotalCount, nil
}

// Get Memberships for User
//
// Returns all of the group memberships for a given user. Note this is only available to group admins. To retrieve group memberships for the user making the API request, use the users/me/memberships endpoint.
// https://developer.box.com/reference#get-all-group-memberships-for-a-user
func (m *Membership) GetMembershipForUserReq(userId string, offset int32, limit int32) *Request {
	var url string
	if limit > 1000 {
		limit = 1000
	}
	url = fmt.Sprintf("%s%s%s%s?offset=%d&limit=%d", m.apiInfo.api.BaseURL, "users/", userId, "/memberships", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}

// Get Memberships for User
//
// Returns all of the group memberships for a given user. Note this is only available to group admins. To retrieve group memberships for the user making the API request, use the users/me/memberships endpoint.
// https://developer.box.com/reference#get-all-group-memberships-for-a-user
func (m *Membership) GetMembershipForUser(userId string, offset int32, limit int32) (outMembership []*Membership, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetMembershipForUserReq(userId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, 0, 0, 0, newApiStatusError(resp.Body)
	}

	memberships := struct {
		TotalCount int           `json:"total_count"`
		Entries    []*Membership `json:"entries"`
		Offset     int           `json:"offset"`
		Limit      int           `json:"limit"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &memberships)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	for _, v := range memberships.Entries {
		v.apiInfo = m.apiInfo
	}
	return memberships.Entries, memberships.Offset, memberships.Limit, memberships.TotalCount, nil
}

// Get Collaborations for Group
//
// Returns all of the group collaborations for a given group. Note this is only available to group admins.
// https://developer.box.com/reference#get-all-collaborations-for-a-group
func (m *Membership) GetCollaborationsForGroupReq(groupId string, offset int32, limit int32) *Request {
	var url string
	if limit > 1000 {
		limit = 1000
	}
	url = fmt.Sprintf("%s%s%s%s?offset=%d&limit=%d", m.apiInfo.api.BaseURL, "groups/", groupId, "/collaborations", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}

// Get Collaborations for Group
//
// Returns all of the group collaborations for a given group. Note this is only available to group admins.
// https://developer.box.com/reference#get-all-collaborations-for-a-group
func (m *Membership) GetCollaborationsForGroup(groupId string, offset int32, limit int32) (outCollaborations []*Collaboration, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetCollaborationsForGroupReq(groupId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, 0, 0, 0, newApiStatusError(resp.Body)
	}

	collabs := struct {
		TotalCount int              `json:"total_count"`
		Entries    []*Collaboration `json:"entries"`
		Offset     int              `json:"offset"`
		Limit      int              `json:"limit"`
	}{}

	err = UnmarshalJSONWrapper(resp.Body, &collabs)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	for _, v := range collabs.Entries {
		v.apiInfo = m.apiInfo
	}
	return collabs.Entries, collabs.Offset, collabs.Limit, collabs.TotalCount, nil
}
