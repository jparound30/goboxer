package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"strings"
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
	return string(*i)
}

//func timeToString(s *time.Time) string {
//	if s == nil {
//		return "<nil>"
//	} else {
//		return s.String()
//	}
//}
//
//func ugToString(s *UserGroupMini) string {
//	if s == nil {
//		return "<nil>"
//	} else {
//		return s.String()
//	}
//}

func ParseResource(jsonEntity []byte) (r BoxResource, err error) {
	decoder := json.NewDecoder(bytes.NewReader(jsonEntity))
	outerStack := 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		var typ string
		switch token {
		case json.Delim('{'):
			outerStack++
			if outerStack != 1 {
				continue
			}
			stack := 0
			var foundTypeField = false
			newDecoder := json.NewDecoder(io.MultiReader(strings.NewReader("{"), decoder.Buffered()))
		InnerLoop:
			for newDecoder.More() {
				token2, err := newDecoder.Token()
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}
				switch token2 {
				case json.Delim('{'):
					stack++
					continue
				case json.Delim('}'):
					stack--
					if stack == 0 {
						break
					}
					continue
				case json.Delim('['), json.Delim(']'):
					continue
				default:
					switch token2.(type) {
					case string:
						if foundTypeField {
							typ = fmt.Sprint(token2)
							break InnerLoop
						}
					default:
						continue
					}
				}
				if token2 == "type" {
					foundTypeField = true
				}
			}
		case json.Delim('}'):
			outerStack--
			continue
		default:
			continue
		}
		dec := json.NewDecoder(bytes.NewReader(jsonEntity))

		switch typ {
		case "folder":
			folder := &Folder{}
			err = dec.Decode(folder)
			r = folder
		case "file":
			file := &File{}
			err = dec.Decode(file)
			r = file
		case "file_version":
			fv := &FileVersion{}
			err = dec.Decode(fv)
			r = fv
		case "user":
			u := &User{}
			err = dec.Decode(u)
			r = u
		case "group":
			g := &Group{}
			err = dec.Decode(g)
			r = g
		case "group_membership":
			gm := &Membership{}
			err = dec.Decode(gm)
			r = gm
		case "collaboration":
			c := &Collaboration{}
			err = dec.Decode(c)
			r = c
		}
	}
	return r, err
}

func UnmarshalJsonWrapper(data []byte, v BoxResource) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		err = xerrors.Errorf("failed to unmarshal response: %w", err)
		return newApiOtherError(err, string(data))
	}
	return nil
}
