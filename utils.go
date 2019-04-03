package goboxer

import "time"

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
func boolToString(b *bool) string {
	if b == nil {
		return "<nil>"
	} else if !*b {
		return "false"
	} else {
		return "true"
	}
}
func timeToString(s *time.Time) string {
	if s == nil {
		return "<nil>"
	} else {
		return s.String()
	}
}

func ugToString(s *UserGroupMini) string {
	if s == nil {
		return "<nil>"
	} else {
		return s.String()
	}
}
