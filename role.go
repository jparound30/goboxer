package goboxer

type Role string

//func (r *Role) MarshalJSON() ([]byte, error) {
//	if r != nil {
//
//	}
//}

func (r *Role) String() string {
	if r == nil {
		return "<nil>"
	}
	return string(*r)
}

func (r *Role) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	} else {
		return []byte(`"` + r.String() + `"`), nil
	}
}

const (
	EDITOR             Role = "editor"
	VIEWER             Role = "viewer"
	PREVIEWER          Role = "previewer"
	UPLOADER           Role = "uploader"
	PREVIEWER_UPLOADER Role = "previewer uploader"
	VIEWER_UPLOADER    Role = "viewer uploader"
	CO_OWNER           Role = "co-owner"
	OWNER              Role = "owner"
)
