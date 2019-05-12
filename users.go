package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
	EnterpriseTypeRolledOut  EnterpriseType = "rolledout"
)

type Enterprise struct {
	Type EnterpriseType
	Id   string
	Name string
}

func (e *Enterprise) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	if e.Type == EnterpriseTypeRolledOut {
		return []byte("null"), nil
	}
	var buffer bytes.Buffer
	buffer.WriteString(`{`)
	buffer.WriteString(`"type":`)
	marshal, err := json.Marshal(e.Type)
	if err != nil {
		return nil, err
	}
	buffer.Write(marshal)
	buffer.WriteString(`,"id":`)
	buffer.WriteString(`"` + e.Id + `"`)
	buffer.WriteString(`,"name":`)
	buffer.WriteString(`"` + e.Name + `"`)
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}

type User struct {
	UserGroupMini
	apiInfo                       *apiInfo
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
	TrackingCodes                 []map[string]string `json:"tracking_codes,omitempty"`
	CanSeeManagedUsers            *bool               `json:"can_see_managed_users,omitempty"`
	IsSyncEnabled                 *bool               `json:"is_sync_enabled,omitempty"`
	IsExternalCollabRestricted    *bool               `json:"is_external_collab_restricted,omitempty"`
	IsExemptFromDeviceLimits      *bool               `json:"is_exempt_from_device_limits,omitempty"`
	IsExemptFromLoginVerification *bool               `json:"is_exempt_from_login_verification,omitempty"`
	Enterprise                    *Enterprise         `json:"enterprise,omitempty"`
	MyTags                        *[]string           `json:"my_tags,omitempty"`
	Hostname                      *string             `json:"hostname,omitempty"`
	IsPlatformAccessOnly          *bool               `json:"is_platform_access_only,omitempty"`
	ExternalAppUserId             *string             `json:"external_app_user_id,omitempty"`
	IsPasswordResetRequired       *bool               `json:"is_password_reset_required,omitempty"`

	NotifyRolledOut *bool `json:"notify,omitempty"`
	changeFlag      uint64
}

const (
	cUserName uint64 = 1 << (iota)
	cUserRole
	cUserLanguage
	cUserIsSyncEnabled
	cUserJobTitle
	cUserPhone
	cUserAddress
	cUserSpaceAmount
	cUserTrackingCodes
	cUserCanSeeMangedUsers
	cUserTimezone
	cUserIsExemptFromDeviceLimits
	cUserIsExemptFromLoginVerification
	cUserIsExternalCollabRestricted
	cUserStatus
	cUserIsPasswordResetRequired
	cUserRollOut
)

func (u *User) ResourceType() BoxResourceType {
	return UserResource
}

func NewUser(api *APIConn) *User {
	return &User{
		apiInfo: &apiInfo{api: api},
	}
}

var UserAllFields = []string{
	"type", "id", "name", "login", "created_at", "modified_at",
	"language", "timezone", "space_amount", "space_used", "max_upload_size",
	"status", "job_title", "phone", "address", "avatar_url", "role",
	"tracking_codes", "can_see_managed_users", "is_sync_enabled",
	"is_external_collab_restricted", "is_exempt_from_device_limits",
	"is_exempt_from_login_verification", "enterprise",
	"my_tags", "hostname", "is_platform_access_only", "external_app_user_id",
}

// Get Current User
//
// Get information about the user who is currently logged in (i.e. the user for whom this access token was generated).
// https://developer.box.com/reference#get-the-current-users-information
func (u *User) GetCurrentUserReq(fields []string) *Request {
	var urlBase string
	var query string

	urlBase = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users/me")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	return NewRequest(u.apiInfo.api, urlBase+query, GET, nil, nil)
}

// Get User
//
// Get information about a user in the enterprise. Requires enterprise administration authorization.
// https://developer.box.com/reference#users
func (u *User) GetCurrentUser(fields []string) (*User, error) {

	req := u.GetCurrentUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Get User
//
// Get information about a user in the enterprise. Requires enterprise administration authorization.
// https://developer.box.com/reference#users
func (u *User) GetUserReq(userId string, fields []string) *Request {
	var urlBase string
	var query string

	urlBase = fmt.Sprintf("%s%s%s", u.apiInfo.api.BaseURL, "users/", userId)
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	return NewRequest(u.apiInfo.api, urlBase+query, GET, nil, nil)
}

// Get User
//
// Get information about a user in the enterprise. Requires enterprise administration authorization.
// https://developer.box.com/reference#users
func (u *User) GetUser(userId string, fields []string) (*User, error) {

	req := u.GetUserReq(userId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO Get User Avatar

// Create User
//
// Create a new managed user in an enterprise. This method only works for Box admins.
// https://developer.box.com/reference#create-an-enterprise-user
func (u *User) CreateUserReq(fields []string) *Request {
	var urlBase string
	var query string
	urlBase = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &User{}
	data.Login = u.Login
	data.Name = u.Name
	if u.changeFlag&cUserRole == cUserRole {
		data.Role = u.Role
	}
	if u.changeFlag&cUserLanguage == cUserLanguage {
		data.Language = u.Language
	}
	if u.changeFlag&cUserIsSyncEnabled == cUserIsSyncEnabled {
		data.IsSyncEnabled = u.IsSyncEnabled
	}
	if u.changeFlag&cUserJobTitle == cUserJobTitle {
		data.JobTitle = u.JobTitle
	}
	if u.changeFlag&cUserPhone == cUserPhone {
		data.Phone = u.Phone
	}
	if u.changeFlag&cUserAddress == cUserAddress {
		data.Address = u.Address
	}
	if u.changeFlag&cUserSpaceAmount == cUserSpaceAmount {
		data.SpaceAmount = u.SpaceAmount
	}
	if u.changeFlag&cUserTrackingCodes == cUserTrackingCodes {
		data.TrackingCodes = u.TrackingCodes
	}
	if u.changeFlag&cUserCanSeeMangedUsers == cUserCanSeeMangedUsers {
		data.CanSeeManagedUsers = u.CanSeeManagedUsers
	}
	if u.changeFlag&cUserTimezone == cUserTimezone {
		data.Timezone = u.Timezone
	}
	if u.changeFlag&cUserIsExemptFromDeviceLimits == cUserIsExemptFromDeviceLimits {
		data.IsExemptFromDeviceLimits = u.IsExemptFromDeviceLimits
	}
	if u.changeFlag&cUserIsExemptFromLoginVerification == cUserIsExemptFromLoginVerification {
		data.IsExemptFromLoginVerification = u.IsExemptFromLoginVerification
	}
	if u.changeFlag&cUserIsExternalCollabRestricted == cUserIsExternalCollabRestricted {
		data.IsExternalCollabRestricted = u.IsExternalCollabRestricted
	}
	if u.changeFlag&cUserStatus == cUserStatus {
		data.Status = u.Status
	}
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(u.apiInfo.api, urlBase+query, POST, nil, bytes.NewReader(bodyBytes))
}

// Create User
//
// Create a new managed user in an enterprise. This method only works for Box admins.
// https://developer.box.com/reference#create-an-enterprise-user
func (u *User) CreateUser(fields []string) (*User, error) {

	req := u.CreateUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u *User) SetLogin(login string) *User {
	u.Login = &login
	return u
}
func (u *User) SetName(name string) *User {
	u.Name = &name
	u.changeFlag |= cUserName
	return u
}
func (u *User) SetRole(role UserRole) *User {
	u.Role = &role
	u.changeFlag |= cUserRole
	return u
}
func (u *User) SetLanguage(language string) *User {
	u.Language = &language
	u.changeFlag |= cUserLanguage
	return u
}
func (u *User) SetIsSyncEnabled(isSyncEnabled bool) *User {
	u.IsSyncEnabled = &isSyncEnabled
	u.changeFlag |= cUserIsSyncEnabled
	return u
}
func (u *User) SetJobTitle(jobTitle string) *User {
	u.JobTitle = &jobTitle
	u.changeFlag |= cUserJobTitle
	return u
}
func (u *User) SetPhone(phone string) *User {
	u.Phone = &phone
	u.changeFlag |= cUserPhone
	return u
}
func (u *User) SetAddress(address string) *User {
	u.Address = &address
	u.changeFlag |= cUserAddress
	return u
}
func (u *User) SetSpaceAmount(spaceAmount int64) *User {
	u.SpaceAmount = spaceAmount
	u.changeFlag |= cUserSpaceAmount
	return u
}
func (u *User) SetTrackingCodes(trackingCodes []map[string]string) *User {
	u.TrackingCodes = trackingCodes
	u.changeFlag |= cUserTrackingCodes
	return u
}
func (u *User) SetCanSeeManagedUsers(canSeeManagedUsers bool) *User {
	u.CanSeeManagedUsers = &canSeeManagedUsers
	u.changeFlag |= cUserCanSeeMangedUsers
	return u
}
func (u *User) SetTimezone(timezone string) *User {
	u.Timezone = &timezone
	u.changeFlag |= cUserTimezone
	return u
}
func (u *User) SetIsExemptFromDeviceLimits(isExemptFromDeviceLimits bool) *User {
	u.IsExemptFromDeviceLimits = &isExemptFromDeviceLimits
	u.changeFlag |= cUserIsExemptFromDeviceLimits
	return u
}

func (u *User) SetIsExemptFromLoginVerification(isExemptFromLoginVerification bool) *User {
	u.IsExemptFromLoginVerification = &isExemptFromLoginVerification
	u.changeFlag |= cUserIsExemptFromLoginVerification
	return u
}
func (u *User) SetIsExternalCollabRestricted(isExternalCollabRestricted bool) *User {
	u.IsExternalCollabRestricted = &isExternalCollabRestricted
	u.changeFlag |= cUserIsExternalCollabRestricted
	return u
}
func (u *User) SetStatus(status UserStatus) *User {
	u.Status = &status
	u.changeFlag |= cUserStatus
	return u
}
func (u *User) SetIsPasswordResetRequired(b bool) *User {
	u.IsPasswordResetRequired = &b
	u.changeFlag |= cUserIsPasswordResetRequired
	return u
}
func (u *User) SetRollOutOfEnterprise(notify bool) *User {
	u.NotifyRolledOut = &notify
	u.Enterprise = &Enterprise{Type: EnterpriseTypeRolledOut}
	u.changeFlag |= cUserRollOut
	return u
}

// Update User
//
// Update the information for a user.
// https://developer.box.com/reference#update-a-users-information
func (u *User) UpdateUserReq(userId string, fields []string) *Request {
	var urlBase string
	var query string
	urlBase = fmt.Sprintf("%s%s%s", u.apiInfo.api.BaseURL, "users/", userId)

	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &User{}
	if u.changeFlag&cUserName == cUserName {
		data.Name = u.Name
	}
	if u.changeFlag&cUserRole == cUserRole {
		data.Role = u.Role
	}
	if u.changeFlag&cUserLanguage == cUserLanguage {
		data.Language = u.Language
	}
	if u.changeFlag&cUserIsSyncEnabled == cUserIsSyncEnabled {
		data.IsSyncEnabled = u.IsSyncEnabled
	}
	if u.changeFlag&cUserJobTitle == cUserJobTitle {
		data.JobTitle = u.JobTitle
	}
	if u.changeFlag&cUserPhone == cUserPhone {
		data.Phone = u.Phone
	}
	if u.changeFlag&cUserAddress == cUserAddress {
		data.Address = u.Address
	}
	if u.changeFlag&cUserSpaceAmount == cUserSpaceAmount {
		data.SpaceAmount = u.SpaceAmount
	}
	if u.changeFlag&cUserTrackingCodes == cUserTrackingCodes {
		data.TrackingCodes = u.TrackingCodes
	}
	if u.changeFlag&cUserCanSeeMangedUsers == cUserCanSeeMangedUsers {
		data.CanSeeManagedUsers = u.CanSeeManagedUsers
	}
	if u.changeFlag&cUserTimezone == cUserTimezone {
		data.Timezone = u.Timezone
	}
	if u.changeFlag&cUserIsExemptFromDeviceLimits == cUserIsExemptFromDeviceLimits {
		data.IsExemptFromDeviceLimits = u.IsExemptFromDeviceLimits
	}
	if u.changeFlag&cUserIsExemptFromLoginVerification == cUserIsExemptFromLoginVerification {
		data.IsExemptFromLoginVerification = u.IsExemptFromLoginVerification
	}
	if u.changeFlag&cUserIsExternalCollabRestricted == cUserIsExternalCollabRestricted {
		data.IsExternalCollabRestricted = u.IsExternalCollabRestricted
	}
	if u.changeFlag&cUserStatus == cUserStatus {
		data.Status = u.Status
	}
	if u.changeFlag&cUserIsPasswordResetRequired == cUserIsPasswordResetRequired {
		data.IsPasswordResetRequired = u.IsPasswordResetRequired
	}
	if u.changeFlag&cUserRollOut == cUserRollOut {
		data.Enterprise = &Enterprise{Type: EnterpriseTypeRolledOut}
		data.NotifyRolledOut = u.NotifyRolledOut
	}

	bodyBytes, _ := json.Marshal(data)

	return NewRequest(u.apiInfo.api, urlBase+query, PUT, nil, bytes.NewReader(bodyBytes))
}

// Update User
//
// Update the information for a user.
// https://developer.box.com/reference#update-a-users-information
func (u *User) UpdateUser(userId string, fields []string) (*User, error) {

	req := u.UpdateUserReq(userId, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Create App User
//
// Create a new app user in an enterprise.
// https://developer.box.com/reference#create-app-user
func (u *User) CreateAppUserReq(fields []string) *Request {
	var urlBase string
	var query string

	urlBase = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users")
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query = fmt.Sprintf("?%s", fieldsParams)
	}

	data := &User{}
	if u.changeFlag&cUserName == cUserName {
		data.Name = u.Name
	}
	if u.changeFlag&cUserLanguage == cUserLanguage {
		data.Language = u.Language
	}
	if u.changeFlag&cUserJobTitle == cUserJobTitle {
		data.JobTitle = u.JobTitle
	}
	if u.changeFlag&cUserPhone == cUserPhone {
		data.Phone = u.Phone
	}
	if u.changeFlag&cUserAddress == cUserAddress {
		data.Address = u.Address
	}
	if u.changeFlag&cUserSpaceAmount == cUserSpaceAmount {
		data.SpaceAmount = u.SpaceAmount
	}
	if u.changeFlag&cUserCanSeeMangedUsers == cUserCanSeeMangedUsers {
		data.CanSeeManagedUsers = u.CanSeeManagedUsers
	}
	if u.changeFlag&cUserTimezone == cUserTimezone {
		data.Timezone = u.Timezone
	}
	if u.changeFlag&cUserIsExternalCollabRestricted == cUserIsExternalCollabRestricted {
		data.IsExternalCollabRestricted = u.IsExternalCollabRestricted
	}
	if u.changeFlag&cUserStatus == cUserStatus {
		data.Status = u.Status
	}

	b := true
	data.IsPlatformAccessOnly = &b
	bodyBytes, _ := json.Marshal(data)

	return NewRequest(u.apiInfo.api, urlBase+query, POST, nil, bytes.NewReader(bodyBytes))
}

// Create App User
//
// Create a new app user in an enterprise.
// https://developer.box.com/reference#create-app-user
func (u *User) CreateAppUser(fields []string) (*User, error) {

	req := u.CreateAppUserReq(fields)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != http.StatusCreated {
		return nil, newApiStatusError(resp.Body)
	}

	r := &User{apiInfo: &apiInfo{api: u.apiInfo.api}}
	err = UnmarshalJSONWrapper(resp.Body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Delete User
//
// Delete a user.
// https://developer.box.com/reference#delete-an-enterprise-user
func (u *User) DeleteUserReq(userId string, notify bool, force bool) *Request {
	var urlBase string
	urlBase = fmt.Sprintf("%s%s%s?notify=%t&force=%t", u.apiInfo.api.BaseURL, "users/", userId, notify, force)

	return NewRequest(u.apiInfo.api, urlBase, DELETE, nil, nil)
}

// Delete User
//
// Delete a user.
// https://developer.box.com/reference#delete-an-enterprise-user
func (u *User) DeleteUser(userId string, notify bool, force bool) error {

	req := u.DeleteUserReq(userId, notify, force)
	resp, err := req.Send()
	if err != nil {
		return err
	}

	if resp.ResponseCode != http.StatusNoContent {
		return newApiStatusError(resp.Body)
	}

	return nil
}

func (u *User) GetEnterpriseUsersReq(filterTerm string, offset int, limit int, fields []string) *Request {
	var urlBase string
	var query string

	urlBase = fmt.Sprintf("%s%s", u.apiInfo.api.BaseURL, "users")

	if limit > 1000 {
		limit = 1000
	}
	query += fmt.Sprintf("?offset=%d&limit=%d", offset, limit)
	if filterTerm != "" {
		query += fmt.Sprintf("&filter_term=%s", url.QueryEscape(filterTerm))
	}
	if fieldsParams := BuildFieldsQueryParams(fields); fieldsParams != "" {
		query += fmt.Sprintf("&%s", fieldsParams)
	}

	return NewRequest(u.apiInfo.api, urlBase+query, GET, nil, nil)
}
func (u *User) GetEnterpriseUsers(filterTerm string, offset int, limit int, fields []string) (outUsers []*User, outOffset int, outLimit int, outTotalCount int, err error) {

	req := u.GetEnterpriseUsersReq(filterTerm, offset, limit, fields)
	resp, err := req.Send()
	if err != nil {
		return nil, 0, 0, 0, err
	}

	if resp.ResponseCode != http.StatusOK {
		return nil, 0, 0, 0, newApiStatusError(resp.Body)
	}
	users := struct {
		TotalCount int     `json:"total_count"`
		Entries    []*User `json:"entries"`
		Offset     int     `json:"offset"`
		Limit      int     `json:"limit"`
	}{}
	err = UnmarshalJSONWrapper(resp.Body, &users)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	for _, user := range users.Entries {
		user.apiInfo = &apiInfo{api: u.apiInfo.api}
	}
	return users.Entries, users.Offset, users.Limit, users.TotalCount, nil
}
