package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func NewService(ctx context.Context, credentialsPath string, tokenPath string) (*drive.Service, error) {
	creds, err := os.ReadFile(credentialsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	config, err := google.ConfigFromJSON(creds, drive.DriveFileScope)

	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	f, err := os.Open(tokenPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	defer f.Close()

	token := &oauth2.Token{}
	json.NewDecoder(f).Decode(token)

	client := config.Client(ctx, token)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return srv, nil

}

func UploadFile(srv *drive.Service, path string) (*drive.File, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer f.Close()

	// File 구조체 포인터를 생성
	metadata := &drive.File{
		Name: filepath.Base(path),
		// 부모 폴더 지정 가능
		// MimeType 지정 가능
	}

	uploadFile, err := srv.Files.Create(metadata).Media(f).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadFile, nil
}
