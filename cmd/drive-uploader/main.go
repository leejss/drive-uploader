package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/leejss/drive-uploader/internal/gdrive"
)

// 애플리케이션 설정을 위한 상수들
const (
	credentialsPath = "credentials.json"
	tokenPath       = "token.json"
)

func main() {

	filePath := flag.String("file", "", "Path to the file to upload")
	flag.Parse()

	if *filePath == "" {
		log.Println("Usage: go run ./cmd/drive-uploader -file <path-to-your-file>")
		log.Fatal("Error: -file argument is required")
	}

	log.Printf("Attempting to upload file: %s", *filePath)

	ctx := context.Background()

	// Google Drive 서비스 생성 시도
	srv, err := gdrive.NewService(ctx, credentialsPath, tokenPath)
	if err != nil {
		// 에러가 gdrive.ErrorTokenNotFound 종류인지 확인
		if errors.Is(err, gdrive.ErrorTokenNotFound) {
			log.Printf("Token not found, attempting to generate a new one...")

			// 헬퍼 함수를 사용하여 토큰 발급 절차 진행
			config, configErr := getConfig(credentialsPath)
			if configErr != nil {
				log.Fatalf("Fatal: Unable to load client secret file: %v", configErr)
			}

			token, tokenErr := getTokenFromWeb(ctx, config)
			if tokenErr != nil {
				log.Fatalf("Fatal: Unable to get token from web: %v", tokenErr)
			}

			if saveErr := saveToken(tokenPath, token); saveErr != nil {
				log.Printf("Warning: Unable to cache oauth token: %v", saveErr)
			}

			// 토큰 발급 후 서비스 생성 재시도
			log.Println("Token acquired. Retrying to create drive service...")
			srv, err = gdrive.NewService(ctx, credentialsPath, tokenPath)
			if err != nil {
				log.Fatalf("Fatal: Failed to create drive service even after getting token: %v", err)
			}
		} else {
			// 토큰 문제가 아닌 다른 복구 불가능한 에러
			log.Fatalf("Fatal: An unrecoverable error occurred: %v", err)
		}
	}

	log.Println("Successfully created drive service.")

	// 파일 업로드 함수 호출
	log.Println("Starting file upload...")
	uploadedFile, err := gdrive.UploadFile(srv, *filePath)
	if err != nil {
		log.Fatalf("Fatal: Failed to upload file: %v", err)
	}

	// 결과 로깅
	log.Printf("File uploaded successfully!")
	log.Printf("File Name: %s", uploadedFile.Name)
	log.Printf("File ID: %s", uploadedFile.Id)
	log.Printf("View online: %s", uploadedFile.WebViewLink)
}
