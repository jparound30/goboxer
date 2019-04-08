package goboxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
		}
	}
	return r, err
}
