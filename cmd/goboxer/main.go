package main

import (
	"fmt"
	"github.com/jparound30/goboxer"
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

	apiConn := goboxer.NewApiConnWithRefreshToken(clientId, clientSecret, accessToken, refreshToken)

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
	goboxer.Log = &mainState

	// API Usage Example

	// 1. Get Folder Info.
	folder := goboxer.NewFolder(apiConn)
	folderInfo, err := folder.GetInfo("0", goboxer.FolderAllFields)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Folder Info:\n%+v\n", folderInfo)

	createFolder, err := folder.Create("69069008141", "NEW FOLDER", goboxer.FolderAllFields)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("Created Folder Info:\n%+v\n", createFolder)

	uf := goboxer.NewFolder(apiConn)

	uf.SetName("NEW FOLDER " + time.Now().Format("2006-01-02-150405"))

	uf.SetDescription("DESCRIPTION")

	//shPer := nil
	unsharedAt := time.Now().Add(time.Hour * 24 * 365)
	uf.SetSharedLinkCollaborators(unsharedAt)

	uf.FolderUploadEmail = &goboxer.FolderUploadEmail{}
	uf.FolderUploadEmail.SetAccess(goboxer.FolderUploadEmailAccessCollaborators)

	syncState := "not_synced"
	uf.SyncState = &syncState

	uf.Tags = []string{"testtag1"}

	canNonOwnersInvite := false

	uf.CanNonOwnersInvite = &canNonOwnersInvite

	isCollaborationRestrictedToEnterprise := false
	uf.IsCollaborationRestrictedToEnterprise = &isCollaborationRestrictedToEnterprise

	// PendingCollaboration
	collaboration := goboxer.NewCollaboration(apiConn)
	pendingList, _, _, outTotalCount, err := collaboration.PendingCollaborations(0, 1000, goboxer.CollaborationAllFields)
	if err != nil {
		fmt.Printf("pendingConnection: count=%d\n", outTotalCount)
		for _, v := range pendingList {
			fmt.Println(v)
		}
	}

	if createFolder.ID != nil {
		_, _ = uf.Update(*createFolder.ID, goboxer.FolderAllFields)

		collaboration.SetItem(goboxer.TYPE_FOLDER, *createFolder.ID)
		collaboration.SetCanViewPath(true)
		collaboration.SetRole(goboxer.VIEWER)
		collaboration.SetAccessibleByEmailForNewUser(goboxer.TYPE_USER, "goboxer00001@example.com")
		createCollab, err := collaboration.Create(goboxer.CollaborationAllFields, false)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
		fmt.Printf("created collab: %s\n", createCollab)

		updatedCollab, err := collaboration.Update(*createCollab.ID, goboxer.UPLOADER, nil, nil, goboxer.CollaborationAllFields)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
		fmt.Printf("updated collab: %s\n", updatedCollab)

		_, _ = uf.Copy(*createFolder.ID, "69069008141", "COPY_"+*uf.Name, goboxer.FolderAllFields)

		collabs, _ := uf.Collaborations("69069008141", goboxer.CollaborationAllFields)
		for _, collab := range collabs {
			_, _ = fmt.Printf("%s\n", collab)
			info, err := collab.GetInfo(*collab.ID, goboxer.CollaborationAllFields)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("[COLLAB] %s\n", info)
		}
		_ = uf.Delete(*createFolder.ID, false)
	}

}

type Main struct {
}

func (*Main) RequestDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) ResponseDumpf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Debugf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Infof(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Warnf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Errorf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Fatalf(format string, args ...interface{}) {
	fmt.Printf("[goboxer] "+format, args...)
}

func (*Main) Success(apiConn *goboxer.ApiConn) {
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

func (*Main) Fail(apiConn *goboxer.ApiConn, err error) {
	fmt.Printf("%v\n", err)
}
