package gobox //import "github.com/jparound30/gobox"

const (
	VERSION = "0.0.1"
)

var (
	Log Logger = nil
)

type apiInfo struct {
	api *ApiConn
}
