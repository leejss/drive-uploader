package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	credentialPath string
	tokenPath      string
)

const (
	configDirName = ".drive-uploader"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory: ", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	// Check stat first, then create if not exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			log.Fatal("Failed to create config directory: ", err)
		}
		fmt.Printf("Config directory created: %s\n", configDir)
	}

	credentialPath = filepath.Join(configDir, "credential.json")
	tokenPath = filepath.Join(configDir, "token.json")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: drive-uploader <command> <arguments>")
		fmt.Println(("Available commands: upload, auth"))
		os.Exit(1) // 비정상 종료
	}

	switch os.Args[1] {
	case "upload":
		return
	case "auth":
		handleAuth(os.Args[2:])
	default:
		fmt.Println("Invalid command")
		return
	}

}

func handleAuth(args []string) {
	authCmd := flag.NewFlagSet("auth", flag.ExitOnError)

	authCmd.Parse(args)
	if authCmd.NArg() == 0 {
		fmt.Println("Usage: drive-uploader auth <action>")
		fmt.Println("Available actions: list, login, logout")
		os.Exit(1)
	}

	action := authCmd.Arg(0)

	switch action {
	case "list":
		fmt.Println("현재 인증 상태를 확인합니다...")

		if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
			fmt.Println("\U0001f4e2 상태: 인증되지 않음")
			fmt.Println("\u2139️ 조치: `drive-uploader auth login` 명령으로 인증을 진행하세요.")
			return
		}

		fmt.Printf("\u2705 상태: 인증됨\n")
		fmt.Printf("\U0001f4c1 토큰 파일 위치: %s\n", tokenPath)
		return
	case "login":
		fmt.Println("새로운 Google 계정 인증을 시작합니다...")

		if _, err := os.Stat(tokenPath); err == nil {
			fmt.Print("이미 인증 정보가 있습니다. 덮어쓰시겠습니까? (y/N): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("인증을 취소했습니다.")
				return
			}
		}

		ctx := context.Background()
		config, err := getConfig(credentialPath)
		if err != nil {
			log.Fatalf("인증 설정 오류:\n%v", err)
		}

		token, err := getTokenFromWeb(ctx, config)
		if err != nil {
			log.Fatalf("오류: 웹에서 토큰을 가져올 수 없습니다: %v", err)
		}

		if err := saveToken(tokenPath, token); err != nil {
			log.Fatalf("오류: 토큰을 파일에 저장할 수 없습니다: %v", err)
		}

		fmt.Println("✅ 인증에 성공했으며 token.json 파일에 저장되었습니다.")
		return
	case "logout":
		fmt.Println("인증 정보를 삭제합니다...")

		if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
			fmt.Println("\u2139️ 이미 로그아웃 상태입니다.")
			return
		}

		if err := os.Remove(tokenPath); err != nil {
			log.Fatalf("\u274c 오류: 토큰 파일을 삭제할 수 없습니다: %v", err)
		}

		fmt.Println("\u2705 성공적으로 로그아웃되었습니다.")
		return
	default:
		fmt.Println("Invalid action")
		return
	}

}
