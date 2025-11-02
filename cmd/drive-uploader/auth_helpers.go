package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	stateToken = "state-token"
)

// getConfig은 credentials.json을 읽어 OAuth2 설정을 생성합니다.
func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		// os.IsNotExist 오류인지 확인하여 더 구체적인 안내를 제공
		if os.IsNotExist(err) {
			errorMsg := `
오류: credentials.json 파일을 찾을 수 없습니다.
이 파일은 Google Cloud Console에서 다운로드해야 합니다.

1. 아래 주소로 이동하세요:
   https://console.cloud.google.com/apis/credentials

2. "OAuth 2.0 클라이언트 ID" 섹션에서 해당 클라이언트 이름을 찾으세요.

3. 오른쪽에 있는 다운로드(↓) 아이콘을 클릭하여 JSON 파일을 다운로드하세요.

4. 다운로드한 파일의 이름을 'credentials.json'으로 변경하고,
   이 프로그램과 같은 디렉토리에 위치시켜 주세요.
`
			return nil, fmt.Errorf("%s", errorMsg)
		}
		// 다른 종류의 파일 읽기 오류일 경우
		return nil, fmt.Errorf("credentials.json 파일을 읽는 중 오류 발생: %w", err)
	}
	// drive.DriveFileScope는 파일 생성/수정/삭제 권한을 포함합니다.
	return google.ConfigFromJSON(b, drive.DriveFileScope)
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	// 로컬 서버 포트 설정
	config.RedirectURL = "http://localhost:8090/auth/callback"

	authURL := config.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser and authorize the app:")
	fmt.Printf("%v\n\n", authURL)

	// 채널로 authorization code 수신
	codeChan := make(chan string)
	errChan := make(chan error)

	// 로컬 HTTP 서버 시작
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code in URL")
			http.Error(w, "No authorization code", http.StatusBadRequest)
			return
		}
		codeChan <- code
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "✅ Authorization successful! You can close this window and return to the terminal.")
	})

	server := &http.Server{Addr: ":8090"}
	go server.ListenAndServe()
	defer server.Shutdown(ctx)

	var authCode string
	select {
	case authCode = <-codeChan:
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return config.Exchange(ctx, authCode)
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open token file for writing: %w", err)
	}
	defer f.Close()
	// 파일을 열고 json형태의 텍스트를 write
	return json.NewEncoder(f).Encode(token)
}
