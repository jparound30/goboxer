package goboxer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserGroupMini struct {
	Type  *UserGroupType `json:"type,omitempty"`
	ID    *string        `json:"id,omitempty"`
	Name  *string        `json:"name,omitempty"`
	Login *string        `json:"login,omitempty"`
}

func (u *UserGroupMini) IsUser() bool {
	if u == nil {
		return false
	} else if *u.Type == "group" {
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

type CollaborationStatus string

func (cs *CollaborationStatus) MarshalJSON() ([]byte, error) {
	if cs == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + cs.String() + `"`), nil
	}
}

func (cs *CollaborationStatus) String() string {
	if cs == nil {
		return "<nil>"
	}
	return string(*cs)
}

const (
	COLLABORATION_STATUS_ACCEPTED CollaborationStatus = "accepted"
	COLLABORATION_STATUS_PENDING  CollaborationStatus = "pending"
	COLLABORATION_STATUS_REJECTED CollaborationStatus = "rejected"
)

type Collaboration struct {
	apiInfo        *apiInfo             `json:"-"`
	Type           *string              `json:"type,omitempty"`
	ID             *string              `json:"id,omitempty"`
	CreatedBy      *UserGroupMini       `json:"created_by,omitempty"`
	CreatedAt      *time.Time           `json:"created_at,omitempty"`
	ModifiedAt     *time.Time           `json:"modified_at,omitempty"`
	ExpiresAt      *time.Time           `json:"modified_at,omitempty"`
	Status         *CollaborationStatus `json:"status,omitempty"`
	AccessibleBy   *UserGroupMini       `json:"accessible_by"`
	InviteEmail    *string              `json:"invite_email,omitempty"`
	Role           *Role                `json:"role,omitempty"`
	AcknowledgedAt *time.Time           `json:"modified_at,omitempty"`
	Item           *ItemMini            `json:"item,omitempty"`
	CanViewPath    *bool                `json:"can_view_path,omitempty"`
}

func NewCollaboration(api *ApiConn) *Collaboration {
	return &Collaboration{apiInfo: &apiInfo{api: api}}
}

func (c *Collaboration) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("{ Type:%s, ID:%s, CreatedBy:%s, CreatedAt:%s, Modified:%s, ExpiresAt:%s, Status:%s,"+
		" AccessibleBy:%s, InviteEmail:%s, Role:%s, AcknowledgedAt:%s, Item:%s, CanViewPath:%s }",
		toString(c.Type), toString(c.ID), c.CreatedBy, c.CreatedAt, c.ModifiedAt,
		c.ExpiresAt, c.Status.String(), c.AccessibleBy, toString(c.InviteEmail),
		c.Role.String(), c.AcknowledgedAt, c.Item.String(), boolToString(c.CanViewPath))
}

var CollaborationAllFields = []string{"type", "id", "item", "accessible_by", "role", "expires_at",
	"can_view_path", "status", "acknowledged_at", "created_by",
	"created_at", "modified_at", "invite_email"}

func (c *Collaboration) GetInfoReq(collaborationId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId, BuildFieldsQueryParams(fields))

	return NewRequest(c.apiInfo.api, url, GET, nil, nil)
}

func (c *Collaboration) GetInfo(collaborationId string, fields []string) (*Collaboration, error) {

	req := c.GetInfoReq(collaborationId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get collaboration info"))
		return nil, err
	}

	r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Collaboration) SetItem(typ ItemType, id string) *Collaboration {
	c.Item = &ItemMini{
		ID:   &id,
		Type: &typ,
	}
	return c
}

func (c *Collaboration) SetAccessibleById(typ UserGroupType, id string) *Collaboration {
	c.AccessibleBy = &UserGroupMini{
		Type: &typ,
		ID:   &id,
	}
	return c
}

func (c *Collaboration) SetAccessibleByEmailForNewUser(typ UserGroupType, login string) *Collaboration {
	c.AccessibleBy = &UserGroupMini{
		Type:  &typ,
		Login: &login,
	}
	return c
}
func (c *Collaboration) SetRole(role Role) *Collaboration {
	c.Role = &role
	return c
}
func (c *Collaboration) SetCanViewPath(canViewPath bool) *Collaboration {
	c.CanViewPath = &canViewPath
	return c
}

func (c *Collaboration) CreateReq(fields []string, notify bool) *Request {
	var url string
	url = fmt.Sprintf("%s%s?notify=%t&%s", c.apiInfo.api.BaseURL, "collaborations", notify, BuildFieldsQueryParams(fields))
	bodyBytes, err := json.Marshal(c)
	if err != nil {
		// FIXME error handlingggggg....
		fmt.Println(err)
	}
	return NewRequest(c.apiInfo.api, url, POST, nil, bytes.NewReader(bodyBytes))
}
func (c *Collaboration) Create(fields []string, notify bool) (*Collaboration, error) {
	req := c.CreateReq(fields, notify)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to create collaboration"))
		return nil, err
	}

	r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
	err = json.Unmarshal(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Collaboration) UpdateReq(collaborationId string, role Role, status *CollaborationStatus, canViewPath *bool, fields []string) *Request {
	b := struct {
		Role        Role                 `json:"role"`
		Status      *CollaborationStatus `json:"status,omitempty"`
		CanViewPath *bool                `json:"can_view_path,omitempty"`
	}{
		Role:        role,
		Status:      status,
		CanViewPath: canViewPath,
	}
	var url string
	url = fmt.Sprintf("%s%s%s?%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId, BuildFieldsQueryParams(fields))
	bodyBytes, _ := json.Marshal(b)
	return NewRequest(c.apiInfo.api, url, PUT, nil, bytes.NewReader(bodyBytes))
}

func (c *Collaboration) Update(collaborationId string, role Role, status *CollaborationStatus, canViewPath *bool, fields []string) (*Collaboration, error) {
	req := c.UpdateReq(collaborationId, role, status, canViewPath, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	switch resp.ResponseCode {
	case http.StatusOK:
		r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
		err = json.Unmarshal(resp.Body, r)
		if err != nil {
			return nil, err
		}
		return r, nil
	case http.StatusNoContent:
		// If the role field is changed from co-owner to owner, then the collaboration object is deleted
		// and a new one is created with the previous owner granted access as a co-owner.
		// A 204 response is returned in this case.
		return nil, nil
	default:
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to update collaboration"))
		return nil, err
	}
}

func (c *Collaboration) DeleteReq(collaborationId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId)
	return NewRequest(c.apiInfo.api, url, DELETE, nil, nil)
}

func (c *Collaboration) Delete(collaborationId string) error {
	req := c.DeleteReq(collaborationId)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	switch resp.ResponseCode {
	case http.StatusNoContent:
		return nil
	default:
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to delete collaboration"))
		return err
	}
}

// Get all pending collaboration invites for a user.
func (c *Collaboration) PendingCollaborationsReq(offset int, limit int, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s?status=%s&offset=%d&limit=%d&%s",
		c.apiInfo.api.BaseURL, "collaborations", "pending", offset, limit, BuildFieldsQueryParams(fields))

	return NewRequest(c.apiInfo.api, url, GET, nil, nil)
}

// Get all pending collaboration invites for a user.
func (c *Collaboration) PendingCollaborations(offset int, limit int, fields []string) (pendingList []*Collaboration, outOffset int, outLimit int, outTotalCount int, err error) {
	req := c.PendingCollaborationsReq(offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, offset, limit, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get pending collaborations info"))
		return nil, offset, limit, 0, err
	}

	r := struct {
		TotalCount int              `json:"total_count"`
		Offset     int              `json:"offset"`
		Limit      int              `json:"limit"`
		Entries    []*Collaboration `json:"entries"`
	}{}
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, offset, limit, 0, err
	}

	for _, v := range r.Entries {
		v.apiInfo = c.apiInfo
	}
	return r.Entries, r.Offset, r.Limit, r.TotalCount, nil
}
