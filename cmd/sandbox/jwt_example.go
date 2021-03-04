package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jparound30/goboxer"
)

const (
	rootFolderId      = "0"
	jwtConfigFilePath = "./config.json"
)

func main() {

	var userId string
	flag.StringVar(&userId, "userId", "0", "your userId")
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(f.Name)); s != "" {
			_ = f.Value.Set(s)
		}
	})
	flag.Parse()

	if userId == "0" {
		log.Fatalf("userId must be non-zero")
	}

	log.Printf(`=====
 START REQUEST JWT TOKEN for ENTERPRISE.
=====
`)
	configFile, err := os.Open(jwtConfigFilePath)
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE %s", err)
	}
	loader := goboxer.JwtConfigDefaultLoader{}
	apiConn, err := goboxer.NewAPIConnWithJwtConfig(configFile, loader)
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE: %s", err)
	}
	folder := goboxer.NewFolder(apiConn)
	getInfo, err := folder.GetInfo(rootFolderId, nil)
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE GET FOLDER: %s", err)
	}
	log.Printf("%v\n", getInfo)

	err = apiConn.Refresh()
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE API TOKEN REFRESH: %s", err)
	}
	err = apiConn.Authenticate("")
	if err == nil {
		log.Fatalf("Failed: ENTERPRISE Authenticate with AUTHCODE: %s", err)
	}

	_, _ = configFile.Seek(0, io.SeekStart)

	log.Printf(`=====
 START REQUEST JWT TOKEN for User.
=====
`)
	apiConnUser, err := goboxer.NewAPIConnWithJwtConfigForUser(configFile, loader, userId)
	if err != nil {
		log.Fatalf("Failed: USER")
	}
	folder = goboxer.NewFolder(apiConnUser)
	getInfo, err = folder.GetInfo(rootFolderId, nil)
	if err != nil {
		log.Fatalf("Failed: USER GET FOLDER")
	}
	log.Printf("%v\n", getInfo)

	err = apiConnUser.Refresh()
	if err != nil {
		log.Fatalf("Failed: USER API TOKEN REFRESH: %s", err)
	}
	err = apiConnUser.Authenticate("")
	if err == nil {
		log.Fatalf("Failed: USER Authenticate with AUTHCODE: %s", err)
	}

}
