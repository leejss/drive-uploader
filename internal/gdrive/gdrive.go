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

func UploadFile(srv *drive.Service, path string, parentId string) (*drive.File, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer f.Close()

	// File 구조체 포인터를 생성
	metadata := &drive.File{
		Name: filepath.Base(path),
	}

	// parentId가 제공된 경우 부모 폴더 설정
	if parentId != "" {
		metadata.Parents = []string{parentId}
	}

	uploadFile, err := srv.Files.Create(metadata).Media(f).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadFile, nil
}

func UploadFolder(srv *drive.Service, rootPath string) error {

	info, err := os.Stat(rootPath)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", rootPath)
		}

		return fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", rootPath)
	}

	folderName := filepath.Base(rootPath)
	rootFolder, err := findOrCreateDriveFolder(srv, folderName, "")

	if err != nil {
		return fmt.Errorf("failed to create root folder: %w", err)
	}

	fmt.Printf("Root folder created: %s (ID: %s)\n", rootFolder.Name, rootFolder.Id)

	// 폴더 경로와 Drive 폴더 ID를 매핑하는 맵
	folderMap := make(map[string]string)
	folderMap[rootPath] = rootFolder.Id

	// rooPath -> rootFolder.Id

	// filepath.Walk를 사용하여 디렉토리 순회
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		// 루트 경로는 이미 처리했으므로 스킵
		if path == rootPath {
			return nil
		}

		// 부모 디렉토리 경로 가져오기
		parentPath := filepath.Dir(path)
		parentId, exists := folderMap[parentPath] // check if parentId exists

		if !exists {
			return fmt.Errorf("parent folder not found for path: %s", path)
		}

		if info.IsDir() {
			// 디렉토리인 경우: Drive에 폴더 생성
			folder, err := findOrCreateDriveFolder(srv, info.Name(), parentId)
			if err != nil {
				return fmt.Errorf("failed to create folder %s: %w", info.Name(), err)
			}

			// 폴더 맵에 추가
			folderMap[path] = folder.Id
			fmt.Printf("Folder created: %s (ID: %s)\n", info.Name(), folder.Id)
		} else {
			// 파일인 경우: 업로드
			uploadedFile, err := UploadFile(srv, path, parentId)
			if err != nil {
				return fmt.Errorf("failed to upload file %s: %w", info.Name(), err)
			}

			fmt.Printf("File uploaded: %s (ID: %s)\n", uploadedFile.Name, uploadedFile.Id)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upload folder: %w", err)
	}

	return nil
}

func findOrCreateDriveFolder(srv *drive.Service, folderName string, parentId string) (*drive.File, error) {

	q := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder'", folderName)

	if parentId != "" {
		q += fmt.Sprintf(" and '%s' in parents", parentId)
	}

	files, err := srv.Files.List().Q(q).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	if len(files.Files) > 0 {
		return files.Files[0], nil
	}

	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentId != "" {
		folder.Parents = []string{parentId}
	}

	createdFolder, err := srv.Files.Create(folder).Do()

	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return createdFolder, nil
}
