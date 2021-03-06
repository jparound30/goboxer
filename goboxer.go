package goboxer //import "github.com/jparound30/goboxer"

const (
	VERSION = "0.0.1"
)

var (
	Log Logger = nil
)

type apiInfo struct {
	api *APIConn
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

type BoxResourceType int

const (
	FileResource BoxResourceType = iota + 1
	FileVersionResource
	FolderResource
	WebLinkResource
	UserResource
	GroupResource
	MembershipResource
	CollaborationResource
	CommentResource
	TaskResource
	EventResource
	CollectionResource
)

type BoxResource interface {
	ResourceType() BoxResourceType
}
