package main

import (
	"fmt"
	"github.com/jparound30/goboxer"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	ClientID     = ""
	ClientSecret = ""
	stateStr     = ""
)

func init() {
	rand.Seed(time.Now().UnixNano())
	stateStr = GenerateState()
}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateState() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func tokenGenPageHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		content := `
<html lang="ja">
<head>
</head>
<body>
<form action="/tokenGen" method="post">
	<h1>※対象クライアントIDのRedirectURLが https://localhost と設定されている必要があります</h1>
    <div>
        <label for="client_id">クライアントID</label>
        <input id="client_id" name="client_id" type="text" size="120">
    </div>
    <div>
        <label for="client_secret">クライアント機密コード</label>
        <input id="client_secret" name="client_secret" type="password" size="120">
    </div>
    <button type="submit">トークン取得</button>
</form>
</body>
</html>`
		_, _ = w.Write([]byte(content))
	} else if req.Method == http.MethodPost {
		_ = req.ParseForm()
		// CLIENT_IDとCLIENT_SECRETを受信し、RAM保存
		ClientID = req.PostForm.Get("client_id")
		ClientSecret = req.PostForm.Get("client_secret")

		redirectUrl := fmt.Sprintf("https://account.box.com/api/oauth2/authorize?response_type=code&redirect_uri=https%%3A%%2F%%2Flocalhost&client_id=%s&state=%s", ClientID, stateStr)
		w.Header().Set("Location", redirectUrl)
		w.WriteHeader(302)
	}
}

func writeErrorHtml(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	content := `
<html lang="ja">
<head>
</head>
<body>
%s
</body>
</html>`

	c := fmt.Sprintf(content, err.Error())
	_, _ = w.Write([]byte(c))

}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	// Boxからのリダイレクト受信用
	if req.Method == http.MethodGet {
		_ = req.ParseForm()
		authCode := req.Form.Get("code")
		state := req.Form.Get("state")
		if authCode == "" || state != stateStr {
			http.NotFound(w, req)
			return
		}

		apiConn := goboxer.NewApiConnWithRefreshToken(ClientID, ClientSecret, "", "")
		err := apiConn.Authenticate(authCode)
		if err != nil {
			writeErrorHtml(w, err)
			return
		}
		// 念の為一度リフレッシュを実行
		err = apiConn.Refresh()
		if err != nil {
			writeErrorHtml(w, err)
			return
		}

		newAccessToken := apiConn.AccessToken
		newRefreshToken := apiConn.RefreshToken
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		content := `
<html lang="ja">
<head>
</head>
<body>
<form action="/tokenGen" method="post">
    <div>
        <label for="client_id">クライアントID</label>
        <input id="client_id" name="client_id" type="text" size="120" value="%s">
    </div>
    <div>
        <label for="client_secret">クライアント機密コード</label>
        <input id="client_secret" name="client_secret" type="password" size="120" value="%s">
    </div>
    <div>
        <label for="access_token">アクセストークン</label>
        <input id="access_token" name="access_token" type="text" size="120" value="%s">
    </div>
    <div>
        <label for="refresh_token">リフレッシュトークン</label>
        <input id="refresh_token" name="refresh_token" type="text" size="120" value="%s">
    </div>
    <button type="submit">トークン取得</button>
</form>
</body>
</html>`
		c := fmt.Sprintf(content, ClientID, ClientSecret, newAccessToken, newRefreshToken)
		_, _ = w.Write([]byte(c))
	}
}

// sample
func main() {

	fmt.Printf("Access URL : https://localhost/tokenGen\n")
	http.HandleFunc("/tokenGen", tokenGenPageHandler)
	http.HandleFunc("/", redirectHandler)

	err := http.ListenAndServeTLS(":443", "cert/server.crt", "cert/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
