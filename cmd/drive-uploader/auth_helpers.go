package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// getConfig은 credentials.json을 읽어 OAuth2 설정을 생성합니다.
func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}
	// drive.DriveFileScope는 파일 생성/수정/삭제 권한을 포함합니다.
	return google.ConfigFromJSON(b, drive.DriveFileScope)
}

// getTokenFromWeb은 사용자에게 인증 URL을 보여주고 인증 코드를 받아 토큰을 반환합니다.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser and authorize the app:")
	fmt.Printf("%v\n\n", authURL)
	fmt.Print("Enter the authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	return config.Exchange(ctx, authCode)
}

// saveToken은 발급받은 토큰을 파일에 저장합니다.
func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open token file for writing: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
