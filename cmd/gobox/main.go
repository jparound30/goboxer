package main

import (
	"fmt"
	"github.com/jparound30/gobox"
	"io/ioutil"
	"os"
)

var (
	StateFilename = "apiconnstate.json"
)

// sample
func main() {
	clientId := os.Getenv("_BOX_CL_ID")
	clientSecret := os.Getenv("_BOX_CL_SC")
	accessToken := os.Getenv("_BOX_AT")
	refreshToken := os.Getenv("_BOX_RT")

	apiConn := gobox.NewApiConnWithRefreshToken(clientId, clientSecret, accessToken, refreshToken)

	_, err := os.Stat(StateFilename)
	if err == nil {
		bytes, err := ioutil.ReadFile(StateFilename)
		err = apiConn.RestoreApiConn(bytes)
		if err != nil {
			os.Exit(1)
		}
		err = apiConn.RestoreApiConn(bytes)
		if err != nil {
			os.Exit(1)
		}
	}

	err = apiConn.Refresh()
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
		os.Exit(1)
	}
	fmt.Printf("access_token: %s", apiConn.AccessToken)
	fmt.Printf("refresh_token: %s", apiConn.RefreshToken)

	bytes, err := apiConn.SaveState()
	if err != nil {
		os.Exit(1)
	}
	err = ioutil.WriteFile(StateFilename, bytes, 0666)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
