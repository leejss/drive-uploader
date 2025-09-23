package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	credentialsPath = "credentials.json"
)

func NewService(ctx context.Context, tokenPath string) (*drive.Service, error) {

	// Read credentials from file
	creds, err := os.ReadFile(credentialsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	config, err := google.ConfigFromJSON(creds, drive.DriveFileScope)

	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	client, err := getClient(ctx, config, tokenPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return srv, nil

}

func UploadFile(srv *drive.Service, filePath string) (*drive.File, error) {
	fileToUpload, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileToUpload.Close()

	fileMetadata := &drive.File{
		Name: filepath.Base(filePath),
	}

	file, err := srv.Files.Create(fileMetadata).Media(fileToUpload).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return file, nil

}

func getClient(ctx context.Context, config *oauth2.Config, tokenPath string) (*http.Client, error) {
	token, err := tokenFromFile(tokenPath)
	if err != nil {
		// token = getTokenFromWeb(ctx, config)
		// saveToken(tokenPath, token)
		return nil, fmt.Errorf("token not found at %s: %w", tokenPath, err)
	}

	return config.Client(ctx, token), nil

}

func tokenFromFile(tokenPath string) (*oauth2.Token, error) {
	f, err := os.Open(tokenPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err

}

// func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
// 	// auth url을 생성한다음에 보여준다. 그 다음에 fmt.Scan을 이용하여 입력받는다. 입력받은 값을 가지고 exchange를 호출한다.
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

// 	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
// 	fmt.Print("Enter the authorization code: ")

// 	var authCode string

// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Fatalf("Unable to read authorization code: %v", err)
// 	}

// 	token, err := config.Exchange(ctx, authCode)
// 	if err != nil {
// 		log.Fatalf("Unable to exchange authorization code: %v", err)
// 	}

// 	return token

// }

// func saveToken(tokenPath string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", tokenPath)
// 	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Fatalf("Unable to open token file: %v", err)
// 	}

// 	defer f.Close()
// 	json.NewEncoder(f).Encode(token)
// }
