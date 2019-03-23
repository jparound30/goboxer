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
