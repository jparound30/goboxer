package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jparound30/goboxer"
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

	apiConn := goboxer.NewAPIConnWithRefreshToken(clientId, clientSecret, accessToken, refreshToken)

	_, err := os.Stat(StateFilename)
	if err == nil {
		bytes, err := ioutil.ReadFile(StateFilename)
		err = apiConn.RestoreState(bytes)
		if err != nil {
			os.Exit(1)
		}
		err = apiConn.RestoreState(bytes)
		if err != nil {
			os.Exit(1)
		}
	}

	mainState := Main{}
	apiConn.SetAPIConnRefreshNotifier(&mainState)
	goboxer.Log = &mainState

	// API Usage Example

	start := time.Now()
	fmt.Printf("[START] %s\n", start)

	targetFolderIds := []string{
		"69589266880", "69588242434", "69591157142", "69597259825", "69600176576",
		"69599167725", "69600071414", "69600457044", "69600768131", "69599379966",
		"69603182713", "69603580530", "69601304448", "69603916900", "69604001609",
		"69604286046", "69607429222", "69607068767", "69607160730", "69607582962",
		//"69607502386","69607671519","69607592130","69607778850","69608070238",
		//"69645916308","69652504089","69655213570","69654053686","69652144245",
		//"69649766461","69069157892","69253381354","69305456746","69343120914",
		//"69343596090","69341762595","69373593135","69367558881","69386898921",
		//"69439205746","69439301049","69439582542","69537456882","69536757485",
		//"69537183441","69537799228","69537349413","69537777044","69577769805",
		//"69576200341","69578177798","69588058591","69649741250","69649365904",
		//"69649310195","69649755031",
	}

	batchRequest := goboxer.NewBatchRequest(apiConn)

	var br []*goboxer.Request
	for _, id := range targetFolderIds {
		folder := goboxer.NewFolder(apiConn)
		br = append(br, folder.CollaborationsReq(id, goboxer.CollaborationAllFields))
	}

	for counter := 1; counter < 21; counter++ {
		fmt.Printf("===START c=%d  ===================================\n", counter)
		response, err := batchRequest.ExecuteBatch(br)
		if err != nil {
			return
		}
		if response.ResponseCode == http.StatusOK {
			requests := response.Responses
			raMax := 0
			for i, v := range requests {
				status := v.ResponseCode
				fmt.Printf("\t[%d]: status:%d\n", i, status)
				fmt.Printf("\t[%d]: header:%s\n", i, v.Headers)
				if status == http.StatusTooManyRequests {
					ra, _ := strconv.Atoi(v.Headers.Get(goboxer.HttpHeaderRetryAfter))
					if raMax < ra {
						raMax = ra
					}
				}
			}
			if raMax > 0 {
				time.Sleep(time.Duration(raMax) * time.Second)
			}
		} else if response.ResponseCode == http.StatusTooManyRequests {
			t, _ := strconv.Atoi(response.Headers.Get(goboxer.HttpHeaderRetryAfter))
			time.Sleep(time.Duration(t) * time.Second)
		} else {
			fmt.Printf("ERRORRRRRRRRRRRaasdfsdfasdfasdfasd\t status:%d\n", response.ResponseCode)
		}
		fmt.Printf("===END  =======================================\n")
	}
	end := time.Now()
	fmt.Printf("[END  ] %s\n", end)
	fmt.Printf("elapsed time = %d\n", (end.UnixNano()-start.UnixNano())/1000000)
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
	fmt.Printf(format, args...)
}

func (*Main) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Warnf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (*Main) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (*Main) EnabledLoggingResponseBody() bool {
	return false
}
func (*Main) EnabledLoggingRequestBody() bool {
	return true
}

func (*Main) Success(apiConn *goboxer.APIConn) {
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

func (*Main) Fail(apiConn *goboxer.APIConn, err error) {
	fmt.Printf("%v\n", err)
}
