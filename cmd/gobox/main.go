package main

import (
	"fmt"
	"github.com/jparound30/gobox"
	"os"
)

func main() {
	apiConn := gobox.NewApiConnWithRefreshToken("CLIENT_ID", "CLIENT_SECRET", "ACCESS_TOKEN", "REFRESH_TOKEN")
	err := apiConn.Refresh()
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
		os.Exit(1)
	}
	fmt.Printf("access_token: %s", apiConn.AccessToken)
	fmt.Printf("refresh_token: %s", apiConn.RefreshToken)
}
