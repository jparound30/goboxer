package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
	ExpiresAt      *time.Time           `json:"expires_at,omitempty"`
	Status         *CollaborationStatus `json:"status,omitempty"`
	AccessibleBy   *UserGroupMini       `json:"accessible_by"`
	InviteEmail    *string              `json:"invite_email,omitempty"`
	Role           *Role                `json:"role,omitempty"`
	AcknowledgedAt *time.Time           `json:"acknowledged_at,omitempty"`
	Item           *ItemMini            `json:"item,omitempty"`
	CanViewPath    *bool                `json:"can_view_path,omitempty"`
}

func (c *Collaboration) ResourceType() BoxResourceType {
	return CollaborationResource
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
		toString(c.Type), toString(c.ID), c.CreatedBy, timeToString(c.CreatedAt), timeToString(c.ModifiedAt),
		timeToString(c.ExpiresAt), c.Status.String(), c.AccessibleBy, toString(c.InviteEmail),
		c.Role.String(), timeToString(c.AcknowledgedAt), c.Item.String(), boolToString(c.CanViewPath))
}

var CollaborationAllFields = []string{"type", "id", "item", "accessible_by", "role", "expires_at",
	"can_view_path", "status", "acknowledged_at", "created_by",
	"created_at", "modified_at", "invite_email"}

// Get Collaboration
//
// Get information about a collaboration.
// https://developer.box.com/reference#get-collabs
func (c *Collaboration) GetInfoReq(collaborationId string, fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId)
	if fieldsParam := BuildFieldsQueryParams(fields); fieldsParam != "" {
		query = fmt.Sprintf("?%s", fieldsParam)
	}

	return NewRequest(c.apiInfo.api, url+query, GET, nil, nil)
}

// Get Collaboration
//
// Get information about a collaboration.
// https://developer.box.com/reference#get-collabs
func (c *Collaboration) GetInfo(collaborationId string, fields []string) (*Collaboration, error) {

	req := c.GetInfoReq(collaborationId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Set target Item(file or folder) (for Create)
// TODO Refactoring.　Is this necessary?
func (c *Collaboration) SetItem(typ ItemType, id string) *Collaboration {
	c.Item = &ItemMini{
		ID:   &id,
		Type: &typ,
	}
	return c
}

// Set Accessible for box user (for Create)
// TODO Refactoring.　Is this necessary?
func (c *Collaboration) SetAccessibleById(typ UserGroupType, id string) *Collaboration {
	c.AccessibleBy = &UserGroupMini{
		Type: &typ,
		ID:   &id,
	}
	return c
}

// Set Accessible with email address (for Create)
// TODO Refactoring.　Is this necessary?
func (c *Collaboration) SetAccessibleByEmailForNewUser(login string) *Collaboration {
	typ := TYPE_USER
	c.AccessibleBy = &UserGroupMini{
		Type:  &typ,
		Login: &login,
	}
	return c
}

// Set Role of collaboration (for Create/Update)
// TODO Refactoring.　Is this necessary?
func (c *Collaboration) SetRole(role Role) *Collaboration {
	c.Role = &role
	return c
}

// Set CanViewPath (for Create)
// TODO Refactoring.　Is this necessary?
func (c *Collaboration) SetCanViewPath(canViewPath bool) *Collaboration {
	c.CanViewPath = &canViewPath
	return c
}

// Create Collaboration
//
// Create a new collaboration that grants a user or group access to a file or folder in a specific role.
// https://developer.box.com/reference#add-a-collaboration
func (c *Collaboration) CreateReq(targetItem ItemMini, grantedTo UserGroupMini, role Role, canViewPath *bool, fields []string, notify bool) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s", c.apiInfo.api.BaseURL, "collaborations")

	query = fmt.Sprintf("?notify=%t", notify)
	if fieldsParam := BuildFieldsQueryParams(fields); fieldsParam != "" {
		query = query + "&" + fieldsParam
	}

	body := struct {
		Item         ItemMini      `json:"item"`
		AccessibleBy UserGroupMini `json:"accessible_by"`
		Role         Role          `json:"role"`
		CanViewPath  *bool         `json:"can_view_path,omitempty"`
	}{
		Item:         targetItem,
		AccessibleBy: grantedTo,
		Role:         role,
		CanViewPath:  canViewPath,
	}
	bodyBytes, _ := json.Marshal(body)
	return NewRequest(c.apiInfo.api, url+query, POST, nil, bytes.NewReader(bodyBytes))
}

// Create Collaboration
//
// Create a new collaboration that grants a user or group access to a file or folder in a specific role.
// https://developer.box.com/reference#add-a-collaboration
func (c *Collaboration) Create(targetItem ItemMini, grantedTo UserGroupMini, role Role, canViewPath *bool, fields []string, notify bool) (*Collaboration, error) {
	req := c.CreateReq(targetItem, grantedTo, role, canViewPath, fields, notify)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Update Collaboration
//
// Update a collaboration.
// https://developer.box.com/reference#edit-a-collaboration
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
	var query string

	url = fmt.Sprintf("%s%s%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId)

	if fieldsParam := BuildFieldsQueryParams(fields); fieldsParam != "" {
		query = fmt.Sprintf("?%s", fieldsParam)

	}
	bodyBytes, _ := json.Marshal(b)
	return NewRequest(c.apiInfo.api, url+query, PUT, nil, bytes.NewReader(bodyBytes))
}

// Update Collaboration
//
// Update a collaboration.
// https://developer.box.com/reference#edit-a-collaboration
func (c *Collaboration) Update(collaborationId string, role Role, status *CollaborationStatus, canViewPath *bool, fields []string) (*Collaboration, error) {
	req := c.UpdateReq(collaborationId, role, status, canViewPath, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	switch resp.ResponseCode {
	case http.StatusOK:
		r := &Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
		err = UnmarshalJSONWrapper(resp.Body, r)
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
		return nil, newApiStatusError(resp.Body)
	}
}

// Delete Collaboration
//
// Delete a collaboration.
// https://developer.box.com/reference#edit-a-collaboration
func (c *Collaboration) DeleteReq(collaborationId string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s", c.apiInfo.api.BaseURL, "collaborations/", collaborationId)
	return NewRequest(c.apiInfo.api, url, DELETE, nil, nil)
}

// Delete Collaboration
//
// Delete a collaboration.
// https://developer.box.com/reference#edit-a-collaboration
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
		err = newApiStatusError(resp.Body)
		return err
	}
}

// Pending Collaborations
//
// Get all pending collaboration invites for a user.
// https://developer.box.com/reference#get-pending-collaborations
func (c *Collaboration) PendingCollaborationsReq(offset int, limit int, fields []string) *Request {
	var url string
	var query string

	url = fmt.Sprintf("%s%s", c.apiInfo.api.BaseURL, "collaborations")

	query = fmt.Sprintf("?status=%s&offset=%d&limit=%d", "pending", offset, limit)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = query + fmt.Sprintf("&%s", fieldsParams)
	}

	return NewRequest(c.apiInfo.api, url+query, GET, nil, nil)
}

// Get all pending collaboration invites for a user.
func (c *Collaboration) PendingCollaborations(offset int, limit int, fields []string) (pendingList []*Collaboration, outOffset int, outLimit int, outTotalCount int, err error) {
	req := c.PendingCollaborationsReq(offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, offset, limit, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, offset, limit, 0, newApiStatusError(resp.Body)
	}

	r := struct {
		TotalCount int              `json:"total_count"`
		Offset     int              `json:"offset"`
		Limit      int              `json:"limit"`
		Entries    []*Collaboration `json:"entries"`
	}{}
	err = UnmarshalJSONWrapper(resp.Body, &r)
	if err != nil {
		return nil, offset, limit, 0, err
	}

	for _, v := range r.Entries {
		v.apiInfo = c.apiInfo
	}
	return r.Entries, r.Offset, r.Limit, r.TotalCount, nil
}
