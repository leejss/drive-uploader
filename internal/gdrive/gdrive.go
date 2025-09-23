package gdrive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var ErrorTokenNotFound = errors.New("token file not found")

func NewService(ctx context.Context, credentialsPath string, tokenPath string) (*drive.Service, error) {
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
		return nil, fmt.Errorf("%w", ErrorTokenNotFound)
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
