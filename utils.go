package goboxer

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"slices"
	"strconv"
	"time"
)

const (
	ContentTypeApplicationJson = "application/json"
	ContentTypeFormUrlEncoded  = "application/x-www-form-urlencoded"
)

func BuildFieldsQueryParams(fields []string) string {
	var params = ""
	if fieldsLen := len(fields); fieldsLen != 0 {
		buffer := make([]byte, 0, 512)
		buffer = append(buffer, "fields="...)
		for index, v := range fields {
			buffer = append(buffer, v...)
			if index != fieldsLen-1 {
				buffer = append(buffer, ',')
			}
		}
		params = string(buffer)
	}
	return params
}

func toString(s *string) string {
	if s == nil {
		return "<nil>"
	} else {
		return *s
	}
}
func timeToString(t *time.Time) string {
	if t == nil {
		return "<nil>"
	}
	return t.Format(time.RFC3339)
}
func boolToString(b *bool) string {
	if b == nil {
		return "<nil>"
	} else if !*b {
		return "false"
	} else {
		return "true"
	}
}
func intToString(i *int) string {
	if i == nil {
		return "<nil>"
	}
	return strconv.Itoa(*i)
}

// ParseResource decodes a Box API object into the concrete resource type selected by its
// top-level "type" field. Supported types are file, folder, file_version, user, group,
// group_membership and collaboration.
//
// An error is returned when the JSON is invalid, has no top-level "type" field, or names an
// unsupported type. When the type is supported but decoding the object fails, the partially
// decoded resource is returned together with the error.
func ParseResource(jsonEntity []byte) (BoxResource, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(jsonEntity, &probe); err != nil {
		return nil, err
	}
	r := newResourceForType(probe.Type)
	if r == nil {
		return nil, xerrors.Errorf("failed to parse box resource: unsupported type %q", probe.Type)
	}
	err := json.Unmarshal(jsonEntity, r)
	return r, err
}

// newResourceForType returns an empty resource for the given Box "type" value, or nil if unsupported.
func newResourceForType(typ string) BoxResource {
	switch typ {
	case "folder":
		return &Folder{}
	case "file":
		return &File{}
	case "file_version":
		return &FileVersion{}
	case "user":
		return &User{}
	case "group":
		return &Group{}
	case "group_membership":
		return &Membership{}
	case "collaboration":
		return &Collaboration{}
	default:
		return nil
	}
}

func setApiInfo(r BoxResource, info *apiInfo) {
	switch r.ResourceType() {
	case FolderResource:
		s := r.(*Folder)
		s.apiInfo = info
	case FileResource:
		s := r.(*File)
		s.apiInfo = info
	case FileVersionResource:
		s := r.(*FileVersion)
		s.apiInfo = info
	case UserResource:
		s := r.(*User)
		s.apiInfo = info
	case GroupResource:
		s := r.(*Group)
		s.apiInfo = info
	case MembershipResource:
		s := r.(*Membership)
		s.apiInfo = info
	case CollaborationResource:
		s := r.(*Collaboration)
		s.apiInfo = info
	default:
		// nothing to do
	}
}

func UnmarshalJSONBoxResourceWrapper(data []byte, v BoxResource) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		err = xerrors.Errorf("failed to unmarshal response: %w", err)
		return newApiOtherError(err, string(data))
	}
	return nil
}

func UnmarshalJSONWrapper(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		err = xerrors.Errorf("failed to unmarshal response: %w", err)
		return newApiOtherError(err, string(data))
	}
	return nil
}

// sendAndCheck sends req and verifies that the response status is one of okStatuses.
// Any other status is converted into an ApiStatusError built from the response body.
func sendAndCheck(req *Request, okStatuses ...int) (*Response, error) {
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	if !slices.Contains(okStatuses, resp.ResponseCode) {
		return nil, newApiStatusError(resp.Body)
	}
	return resp, nil
}

// sendAndDecode sends req, verifies the response status with sendAndCheck,
// and decodes the JSON response body into out.
func sendAndDecode(req *Request, out any, okStatuses ...int) (*Response, error) {
	resp, err := sendAndCheck(req, okStatuses...)
	if err != nil {
		return nil, err
	}
	if err := UnmarshalJSONWrapper(resp.Body, out); err != nil {
		return nil, err
	}
	return resp, nil
}
