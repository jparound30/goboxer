package main

import (
	"fmt"
	"github.com/jparound30/gobox"
	"io/ioutil"
	"os"
	"time"
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
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Folder Info:\n%+v", folderInfo)

	createFolder, err := folder.Create("69069008141", "NEW FOLDER", gobox.FolderAllFields)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("Created Folder Info:\n%+v\n", createFolder)

	uf := gobox.NewFolder(apiConn)

	uf.SetName("NEW FOLDER " + time.Now().Format("2006-01-02-150405"))

	uf.SetDescription("DESCRIPTION")

	//shPer := nil
	unsharedAt := time.Now().Add(time.Hour * 24 * 365)
	uf.SetSharedLinkCollaborators(unsharedAt)

	uf.FolderUploadEmail = &gobox.FolderUploadEmail{}
	uf.FolderUploadEmail.SetAccess(gobox.FolderUploadEmailAccessCollaborators)

	syncState := "not_synced"
	uf.SyncState = &syncState

	uf.Tags = []string{"testtag1"}

	canNonOwnersInvite := false

	uf.CanNonOwnersInvite = &canNonOwnersInvite

	isCollaborationRestrictedToEnterprise := false
	uf.IsCollaborationRestrictedToEnterprise = &isCollaborationRestrictedToEnterprise

	if createFolder.ID != nil {
		_, _ = uf.Update(*createFolder.ID, gobox.FolderAllFields)

		_, _ = uf.Copy(*createFolder.ID, "69069008141", "COPY_"+*uf.Name, gobox.FolderAllFields)
		_ = uf.Delete(*createFolder.ID, false)
	}

}

type Main struct {
}

func (*Main) Success(apiConn *gobox.ApiConn) {
	fmt.Printf("access_token: %s\n", apiConn.AccessToken)
	fmt.Printf("refresh_token: %s\n", apiConn.RefreshToken)
	bytes, err := apiConn.SaveState()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	err = ioutil.WriteFile(StateFilename, bytes, 0666)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func (*Main) Fail(apiConn *gobox.ApiConn, err error) {
	fmt.Printf("%v\n", err)
}
