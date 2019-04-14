package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
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
	apiInfo                 *apiInfo                 `json:"-"`
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
func (m *Membership) GetMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships", membershipId)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}
func (m *Membership) GetMembership(membershipId string) (*Membership, error) {

	req := m.GetMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get membership"))
		return nil, err
	}

	membership := Membership{}

	err = json.Unmarshal(resp.Body, &membership)
	if err != nil {
		return nil, err
	}
	membership.apiInfo = m.apiInfo
	return &membership, nil
}

// Create Membership
func (m *Membership) CreateMembershipReq() *Request {
	var url string
	url = fmt.Sprintf("%s%s", m.apiInfo.api.BaseURL, "group_memberships")

	b, err := json.Marshal(m)
	if err != nil {
		// TODO error handling....
		fmt.Println(err)
	}
	return NewRequest(m.apiInfo.api, url, POST, nil, bytes.NewReader(b))
}
func (m *Membership) CreateMembership() (*Membership, error) {

	req := m.CreateMembershipReq()
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to create membership"))
		return nil, err
	}

	membership := Membership{}

	err = json.Unmarshal(resp.Body, &membership)
	if err != nil {
		return nil, err
	}
	membership.apiInfo = m.apiInfo
	return &membership, nil
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
func (m *Membership) UpdateMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships/", membershipId)

	b, err := json.Marshal(m)
	if err != nil {
		// TODO error handling....
		fmt.Println(err)
	}
	return NewRequest(m.apiInfo.api, url, PUT, nil, bytes.NewReader(b))
}
func (m *Membership) UpdateMembership(membershipId string) (*Membership, error) {

	req := m.UpdateMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to update membership"))
		return nil, err
	}

	membership := Membership{}

	err = json.Unmarshal(resp.Body, &membership)
	if err != nil {
		return nil, err
	}
	membership.apiInfo = m.apiInfo
	return &membership, nil
}

// Delete Membership
func (m *Membership) DeleteMembershipReq(membershipId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", m.apiInfo.api.BaseURL, "group_memberships/", membershipId)

	return NewRequest(m.apiInfo.api, url, DELETE, nil, nil)
}
func (m *Membership) DeleteMembership(membershipId string) error {

	req := m.DeleteMembershipReq(membershipId)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to delete membership"))
		return err
	}
	return nil
}

// Get Memberships for Group
func (m *Membership) GetMembershipForGroupReq(groupId string, offset int32, limit int32) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s%s?&offset=%d&limit=%d", m.apiInfo.api.BaseURL, "groups/", groupId, "/memberships", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}
func (m *Membership) GetMembershipForGroup(groupId string, offset int32, limit int32) (outMembership []*Membership, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetMembershipForGroupReq(groupId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get memberships for group"))
		return nil, 0, 0, 0, err
	}

	memberships := struct {
		TotalCount int           `json:"total_count"`
		Entries    []*Membership `json:"entries"`
		Offset     int           `json:"offset"`
		Limit      int           `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &memberships)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return memberships.Entries, memberships.Offset, memberships.Limit, memberships.TotalCount, nil
}

// Get Memberships for User
func (m *Membership) GetMembershipForUserReq(userId string, offset int32, limit int32) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s%s?offset=%d&limit=%d", m.apiInfo.api.BaseURL, "users/", userId, "/memberships", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}
func (m *Membership) GetMembershipForUser(userId string, offset int32, limit int32) (outMembership []*Membership, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetMembershipForUserReq(userId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get memberships for user"))
		return nil, 0, 0, 0, err
	}

	memberships := struct {
		TotalCount int           `json:"total_count"`
		Entries    []*Membership `json:"entries"`
		Offset     int           `json:"offset"`
		Limit      int           `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &memberships)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return memberships.Entries, memberships.Offset, memberships.Limit, memberships.TotalCount, nil
}

// Get Collaborations for Group
func (m *Membership) GetCollaborationsForGroupReq(groupId string, offset int32, limit int32) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s%s?&offset=%d&limit=%d", m.apiInfo.api.BaseURL, "groups/", groupId, "/collaborations", offset, limit)

	return NewRequest(m.apiInfo.api, url, GET, nil, nil)
}
func (m *Membership) GetCollaborationsForGroup(groupId string, offset int32, limit int32) (outCollaborations []*Collaboration, outOffset int, outLimit int, outTotalCount int, err error) {

	req := m.GetCollaborationsForGroupReq(groupId, offset, limit)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get collaborations for group"))
		return nil, 0, 0, 0, err
	}

	collabs := struct {
		TotalCount int              `json:"total_count"`
		Entries    []*Collaboration `json:"entries"`
		Offset     int              `json:"offset"`
		Limit      int              `json:"limit"`
	}{}

	err = json.Unmarshal(resp.Body, &collabs)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return collabs.Entries, collabs.Offset, collabs.Limit, collabs.TotalCount, nil
}
