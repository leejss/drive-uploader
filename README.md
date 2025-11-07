# Drive Uploader

Google Drive에 파일과 폴더를 업로드할 수 있는 커맨드라인 도구입니다.

## 설치

### 소스 코드에서 빌드

```bash
# 리포지토리 클론
git clone https://github.com/leejss/drive-uploader.git
cd drive-uploader

# 의존성 다운로드
go mod download

# 바이너리 빌드
go build -o drive-uploader ./cmd/drive-uploader

# 빌드된 바이너리 실행
./drive-uploader <command> <arguments>
```

### 직접 실행 (개발 환경)

```bash
# 의존성 다운로드
go mod download

# 직접 실행
go run ./cmd/drive-uploader <command> <arguments>
```

## 사용법

### 기본 명령어 구조

```bash
drive-uploader <command> <arguments>
```

### 사용 가능한 명령어

- `auth` - 인증 관리
- `upload` - 파일/폴더 업로드

## 인증 관리

### 인증 상태 확인

```bash
drive-uploader auth list
```

현재 인증 상태를 확인합니다:

- ✅ 인증됨: 토큰 파일이 존재하고 유효한 상태
- 📢 인증되지 않음: 로그인이 필요한 상태

### 로그인

```bash
drive-uploader auth login
```

Google 계정으로 로그인합니다:

1. 브라우저에서 인증 URL이 자동으로 열립니다
2. Google 계정으로 로그인하고 권한을 승인합니다
3. 인증이 완료되면 토큰이 자동으로 저장됩니다

**참고**: 이미 인증 정보가 있는 경우 덮어쓸지 확인합니다.

### 로그아웃

```bash
drive-uploader auth logout
```

저장된 인증 정보를 삭제합니다.

## 파일 업로드

### 단일 파일 업로드

```bash
drive-uploader upload file <파일경로>
```

**예시:**

```bash
# 절대 경로
drive-uploader upload file /path/to/your/file.txt

# 상대 경로
drive-uploader upload file ./documents/report.pdf

# 홈 디렉토리 파일
drive-uploader upload file ~/Downloads/image.jpg
```

### 폴더 업로드

```bash
drive-uploader upload folder <폴더경로>
```

**예시:**

```bash
# 현재 디렉토리의 폴더
drive-uploader upload folder ./my-folder

# 절대 경로의 폴더
drive-uploader upload folder /path/to/folder

# 홈 디렉토리의 폴더
drive-uploader upload folder ~/Documents/project
```

**특징:**

- 폴더 구조를 그대로 유지하여 업로드
- 하위 폴더와 파일을 재귀적으로 처리
- 업로드 진행 상황을 실시간으로 표시

## 설정

### 설정 파일 위치

프로그램은 다음 위치에 설정 파일을 저장합니다:

```
~/.drive-uploader/
├── credential.json    # Google API 인증 정보 (수동 배치 필요)
└── token.json        # OAuth 토큰 (자동 생성)
```

### Google Cloud 설정

1. [Google Cloud Console](https://console.cloud.google.com/)에서 프로젝트 생성
2. Drive API 활성화
3. OAuth 2.0 클라이언트 ID 생성
4. `credentials.json` 파일 다운로드
5. `~/.drive-uploader/credentials.json`로 파일 배치

자세한 설정 방법은 [설치 및 설정 가이드](docs/03_설치_및_설정.md)를 참조하세요.

## 에러 해결

### 인증 관련 오류

```bash
# 인증 상태 확인
drive-uploader auth list

# 재인증
drive-uploader auth logout
drive-uploader auth login
```

### credentials.json 파일 없음

```
오류: credentials.json 파일을 찾을 수 없습니다.
```

**해결방법:**

1. Google Cloud Console에서 OAuth 2.0 클라이언트 ID의 JSON 파일 다운로드
2. 파일명을 `credentials.json`으로 변경
3. `~/.drive-uploader/` 디렉토리에 배치

## 예시 사용 시나리오

### 1. 처음 사용하는 경우

```bash
# 1. 인증 상태 확인
drive-uploader auth list

# 2. 로그인 (처음 사용 시)
drive-uploader auth login

# 3. 파일 업로드
drive-uploader upload file ./test.txt
```

### 2. 프로젝트 폴더 전체 업로드

```bash
# 프로젝트 폴더 업로드
drive-uploader upload folder ./my-project

# 업로드 결과 확인
# ✅ 폴더 업로드 성공
```

### 3. 여러 파일 순차 업로드

```bash
# 각 파일을 개별적으로 업로드
drive-uploader upload file ./doc1.pdf
drive-uploader upload file ./doc2.pdf
drive-uploader upload file ./image.jpg
```
