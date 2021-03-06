package main

import (
	"code.google.com/p/google-api-go-client/drive/v2"
	"code.google.com/p/goauth2/oauth"
    "golangcafe/goauth2sample/authorize"
	"bytes"
	"flag"
	"fmt"
	"log"
)

var (
	cachefile = "cache.json"

	scope = "https://www.googleapis.com/auth/drive"
    // request_urlは使用するAPIのURLを指定して下さい。（この例ではCalendarList）
	request_url = "https://www.googleapis.com/calendar/v3/users/me/calendarList"
    request_token_url = "https://accounts.google.com/o/oauth2/auth"
    auth_token_url = "https://accounts.google.com/o/oauth2/token"
)

func main() {
    flag.Parse()

    // 認証情報の取得（何もなければ、入力を促します）
    // clientID、secret、redirect_urlはDevelopers ConsoleのCredentialsからコピー＆ペーストして下さい。
    var auth authorize.Auth
    var err error
    if auth, err = authorize.GetAuthInfo(); err != nil {
        log.Fatalln("GetAuthInfo: ", err)
    }

    fmt.Println("Start Execute API")

    // 認証コードを引数で受け取る。
    code := flag.Arg(0)

    config := &oauth.Config{
            ClientId:     auth.ClientID,
            ClientSecret: auth.Secret,
            RedirectURL:  auth.RedirectUrl,
            Scope:        scope,
            AuthURL:      request_token_url,
            TokenURL:     auth_token_url,
            TokenCache:   oauth.CacheFile(cachefile),
    }

    transport := &oauth.Transport{Config: config}

    // キャッシュからトークンファイルを取得
    _, err = config.TokenCache.Token()
    if err != nil {
        // キャッシュなし

        // 認証コードなし＝＞ブラウザで認証させるためにURLを出力
        if code == "" {
            url := config.AuthCodeURL("")
            fmt.Println("ブラウザで以下のURLにアクセスし、認証して下さい。")
            log.Fatalln(url)
        }

        // 認証トークンを取得する。（取得後、キャッシュへ）
        _, err = transport.Exchange(code)
        if err != nil {
            log.Fatalln("Exchange: ", err)
        }

    }
	
	// Google DriveのClientを作成
	svc, err := drive.New(transport.Client())
	if err != nil {
		log.Fatalln("drive.New: ", err)
		return
	}

	f := &drive.File {
		Title: "Sample Document",
		Description: "Sample Document by Go",
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte("This document written by Go Program."))
	r, err := svc.Files.Insert(f).Media(buf).Do()
	if err != nil {
		log.Fatalln("drive.New: ", err)		
	}

	fmt.Printf("Created: ID = %v, Title = %v", r.Id, r.Title)
}