package goboxer //import "github.com/jparound30/goboxer"

const (
	VERSION = "0.0.1"
)

var (
	Log Logger = nil
)

type apiInfo struct {
	api *ApiConn
}

type ItemType string

func (i *ItemType) String() string {
	if i == nil {
		return "<nil>"
	}
	return string(*i)
}
func (i *ItemType) MarshalJSON() ([]byte, error) {
	if i == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + i.String() + `"`), nil
	}
}

const (
	TYPE_FILE   ItemType = "file"
	TYPE_FOLDER ItemType = "folder"
)

type UserGroupType string

func (u *UserGroupType) String() string {
	if u == nil {
		return "<nil>"
	}
	return string(*u)
}
func (u *UserGroupType) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + u.String() + `"`), nil
	}
}

const (
	TYPE_USER  UserGroupType = "user"
	TYPE_GROUP UserGroupType = "group"
)

type BoxResource interface {
	Type() string
}
