package gobox

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UserGroupMini struct {
	Type  *string `json:"type,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Login *string `json:"login,omitempty"`
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
	toString := func(s *string) string {
		if s == nil {
			return "<nil>"
		} else {
			return *s
		}
	}
	return fmt.Sprintf("{Type:%s, ID:%s, Name:%s, Login:%s}", toString(u.Type), toString(u.ID), toString(u.Name), toString(u.Login))
}

type Collaboration struct {
	apiInfo        *apiInfo       `json:"-"`
	Type           *string        `json:"type,omitempty"`
	ID             *string        `json:"id,omitempty"`
	CreatedBy      *UserGroupMini `json:"created_by"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	ModifiedAt     *time.Time     `json:"modified_at,omitempty"`
	ExpiresAt      *time.Time     `json:"modified_at,omitempty"`
	Status         *string        `json:"status,omitempty"`
	AccessibleBy   *UserGroupMini `json:"accessible_by"`
	InviteEmail    *string        `json:"invite_email"`
	Role           *string        `json:"role,omitempty"`
	AcknowledgedAt *time.Time     `json:"modified_at,omitempty"`
	Item           *ItemMini      `json:"item,omitempty"`
	CanViewPath    *bool          `json:"can_view_path,omitempty"`
}

func NewCollaboration(api *ApiConn) *Collaboration {
	return &Collaboration{apiInfo: &apiInfo{api: api}}
}

func (c *Collaboration) String() string {
	if c == nil {
		return "<nil>"
	}
	toString := func(s *string) string {
		if s == nil {
			return "<nil>"
		} else {
			return *s
		}
	}
	boolToString := func(s *bool) string {
		if s == nil {
			return "<nil>"
		} else if !*s {
			return "false"
		} else {
			return "true"
		}
	}
	timeToString := func(s *time.Time) string {
		if s == nil {
			return "<nil>"
		} else {
			return s.String()
		}
	}
	ugToString := func(s *UserGroupMini) string {
		if s == nil {
			return "<nil>"
		} else {
			return s.String()
		}
	}
	return fmt.Sprintf("{Type:%s, ID:%s, CreatedBy:%s, CreatedAt:%s, Modified:%s, ExpiresAt:%s, Status:%s,"+
		" AccessibleBy:%s, InviteEmail:%s, Role:%s, AcknowledgedAt:%s, Item:%s, CanViewPath:%s}",
		toString(c.Type), toString(c.ID), ugToString(c.CreatedBy), timeToString(c.CreatedAt), timeToString(c.ModifiedAt),
		timeToString(c.ExpiresAt), toString(c.Status), ugToString(c.AccessibleBy), toString(c.InviteEmail),
		toString(c.Role), timeToString(c.AcknowledgedAt), c.Item.String(), boolToString(c.CanViewPath))
}

var CollaborationAllFields = []string{"type", "id", "item", "accessible_by", "role", "expires_at",
	"can_view_path", "status", "acknowledged_at", "created_by",
	"created_at", "modified_at", "invite_email"}

func (c *Collaboration) GetInfoReq(collabId string, fields []string) *Request {
	var url string
	url = fmt.Sprintf("%s%s%s?%s", c.apiInfo.api.BaseURL, "collaborations/", collabId, BuildFieldsQueryParams(fields))

	return NewRequest(c.apiInfo.api, url, GET, nil, nil)
}

func (c *Collaboration) GetInfo(collabId string, fields []string) (*Collaboration, error) {

	req := c.GetInfoReq(collabId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get collaboratin info"))
		return nil, err
	}

	r := Collaboration{apiInfo: &apiInfo{api: c.apiInfo.api}}
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Collaboration) Create(fields []string, notify bool) (*Collaboration, error) {
	return nil, nil
}

func (c *Collaboration) Update(collabId string, fields []string) (*Collaboration, error) {
	return nil, nil
}
func (c *Collaboration) Delete(collabId string) error {
	return nil
}

// Get all pending collaboration invites for a user.
func (c *Collaboration) PendingCollaborations(offset int, limit int, fields []string) (pendingList []Collaboration, outOffset int, outLimit int, outTotalCount int, err error) {
	var url string
	url = fmt.Sprintf("%s%s?status=%s&offset=%d&limit=%d&%s",
		c.apiInfo.api.BaseURL, "collaborations", "pending", offset, limit, BuildFieldsQueryParams(fields))

	req := NewRequest(c.apiInfo.api, url, GET, nil, nil)
	resp, err := req.Send()
	if err != nil {
		return nil, offset, limit, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		// TODO improve error handling...
		err = errors.New(fmt.Sprintf("faild to get pending collaboratins info"))
		return nil, offset, limit, 0, err
	}

	r := struct {
		TotalCount int             `json:"total_count"`
		Offset     int             `json:"offset"`
		Limit      int             `json:"limit"`
		Entries    []Collaboration `json:"entries"`
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
