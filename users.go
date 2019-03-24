package goboxer

import "time"

type UserStatus string

const (
	UserStatusActive                 UserStatus = "active"
	UserStatusInactive               UserStatus = "inactive"
	UserStatusCannotDeleteEdit       UserStatus = "cannot_delete_edit"
	UserStatusCannotDeleteEditUpload UserStatus = "cannot_delete_edit_upload"
)

func (us *UserStatus) String() string {
	if us == nil {
		return "<nil>"
	}
	return string(*us)
}

func (us *UserStatus) MarshalJSON() ([]byte, error) {
	if us == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + us.String() + `"`), nil
	}
}

type UserRole string

func (ur *UserRole) String() string {
	if ur == nil {
		return "<nil>"
	}
	return string(*ur)
}

func (ur *UserRole) MarshalJSON() ([]byte, error) {
	if ur == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + ur.String() + `"`), nil
	}
}

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleCoAdmin UserRole = "coadmin"
	UserRoleUser    UserRole = "user"
)

type EnterpriseType string

func (et *EnterpriseType) String() string {
	if et == nil {
		return "<nil>"
	}
	return string(*et)
}

func (et *EnterpriseType) MarshalJSON() ([]byte, error) {
	if et == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + et.String() + `"`), nil
	}
}

const (
	EnterpriseTypeEnterprise EnterpriseType = "enterprise"
	EnterpriseTypeUser       EnterpriseType = "user"
)

type Enterprise struct {
	Type EnterpriseType
	Id   string
	Name string
}

type User struct {
	UserGroupMini
	CreatedAt                     *time.Time          `json:"created_at,omitempty"`
	ModifiedAt                    *time.Time          `json:"modified_at,omitempty"`
	Language                      *string             `json:"language,omitempty"`
	Timezone                      *string             `json:"timezone,omitempty"`
	SpaceAmount                   int64               `json:"space_amount,omitempty"`
	SpaceUsed                     int64               `json:"space_used,omitempty"`
	MaxUploadSize                 int                 `json:"max_upload_size,omitempty"`
	Status                        *UserStatus         `json:"status,omitempty"`
	JobTitle                      *string             `json:"job_title,omitempty"`
	Phone                         *string             `json:"phone,omitempty"`
	Address                       *string             `json:"address,omitempty"`
	AvatarUrl                     *string             `json:"avatar_url,omitempty"`
	Role                          *UserRole           `json:"role,omitempty"`
	TrackingCodes                 []map[string]string `json:"tracking_codes"`
	CanSeeManagedUsers            *bool               `json:"can_see_managed_users,omitempty"`
	IsSyncEnabled                 *bool               `json:"is_sync_enabled,omitempty"`
	IsExternalCollabRestricted    *bool               `json:"is_external_collab_restricted,omitempty"`
	IsExemptFromDeviceLimits      *bool               `json:"is_exempt_from_device_limits,omitempty"`
	IsExemptFromLoginVerification *bool               `json:"is_exempt_from_login_verification,omitempty"`
	Enterprise                    *Enterprise         `json:"enterprise,omitempty"`
	MyTags                        *[]string           `json:"my_tags,omitempty"`
	Hostname                      *string             `json:"hostname,omitempty"`
	IsPlatformAccessOnly          *bool               `json:"is_platform_access_only"`
}
