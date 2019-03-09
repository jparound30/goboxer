package gobox

import (
	"time"
)

type UserOrGroup struct {
	Type  *string `json:"type,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Login *string `json:"login,omitempty"`
}

type Collaboration struct {
	ApiInfo        *apiInfo     `json:"-"`
	Type           *string      `json:"type,omitempty"`
	ID             *string      `json:"id,omitempty"`
	CreatedBy      *UserOrGroup `json:"created_by"`
	CreatedAt      *time.Time   `json:"created_at,omitempty"`
	ModifiedAt     *time.Time   `json:"modified_at,omitempty"`
	ExpiresAt      *time.Time   `json:"modified_at,omitempty"`
	Status         *string      `json:"status,omitempty"`
	AccessibleBy   *UserOrGroup `json:"accessible_by"`
	InviteEmail    *string      `json:"invite_email"`
	Role           *string      `json:"role,omitempty"`
	AcknowledgedAt *time.Time   `json:"modified_at,omitempty"`
	Item           *ItemMini    `json:"item,omitempty"`
	CanViewPath    *bool        `json:"can_view_path,omitempty"`
}

var CollaborationAllFields = []string{"type", "id", "item", "accessible_by", "role", "expires_at",
	"can_view_path", "status", "acknowledged_at", "created_by",
	"created_at", "modified_at", "invite_email"}
