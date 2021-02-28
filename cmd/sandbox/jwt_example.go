package main

import (
	"flag"
	"fmt"
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

	fmt.Printf(`=====
 START REQUEST JWT TOKEN for ENTERPRISE.
=====
`)
	apiConn, err := goboxer.NewAPIConnWithJwtConfig(jwtConfigFilePath)
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE")
	}
	folder := goboxer.NewFolder(apiConn)
	getInfo, err := folder.GetInfo(rootFolderId, nil)
	if err != nil {
		log.Fatalf("Failed: ENTERPRISE GET FOLDER")
	}
	log.Printf("%v\n", getInfo)

	fmt.Printf(`=====
 START REQUEST JWT TOKEN for User.
=====
`)
	apiConnUser, err := goboxer.NewAPIConnWithJwtConfigForUser(jwtConfigFilePath, userId)
	if err != nil {
		log.Fatalf("Failed: USER")
	}
	folder = goboxer.NewFolder(apiConnUser)
	getInfo, err = folder.GetInfo(rootFolderId, nil)
	if err != nil {
		log.Fatalf("Failed: USER GET FOLDER")
	}
	log.Printf("%v\n", getInfo)

}
