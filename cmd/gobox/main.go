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

	mainState := Main{}
	apiConn.SetApiConnRefreshNotifier(&mainState)

	// API Usage Example

	// 1. Get Folder Info.
	folder := gobox.NewFolder(apiConn)
	folderInfo, err := folder.GetInfo("0", gobox.FolderAllFields)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("Folder Info:\n%+v", folderInfo)
}

type Main struct {
}

func (*Main) Success(apiConn *gobox.ApiConn) {
	fmt.Printf("access_token: %s", apiConn.AccessToken)
	fmt.Printf("refresh_token: %s", apiConn.RefreshToken)
	bytes, err := apiConn.SaveState()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	err = ioutil.WriteFile(StateFilename, bytes, 0666)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}

func (*Main) Fail(apiConn *gobox.ApiConn, err error) {
	fmt.Printf("%v", err)
}
