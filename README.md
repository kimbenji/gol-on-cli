# gol-on-cli

Conway's Game of Life를 터미널에서 실행하는 Go 기반 CLI 애플리케이션입니다.

## 주요 기능

- Conway's Game of Life 시뮬레이션 실행
- CLI 옵션(`--help`, `--version`, `--fps`, `--seed`, `--pattern-url`) 지원
- 외부 패턴 URL 로딩 지원

## 로컬에서 실행

### 요구사항

- Go 1.25.1+

### 빌드

```bash
go build -o gol-on-cli ./cmd/gol-on-cli
```

### 실행

```bash
./gol-on-cli --help
./gol-on-cli --fps 15
```

## GitHub Actions로 빌드/릴리스

이 저장소에는 `.github/workflows/ci-release.yml` 워크플로우가 포함되어 있습니다.

- Pull Request / `main` 푸시 시: `go test ./...` 실행
- `v*` 태그 푸시 시:
  - Linux/macOS/Windows용 바이너리 빌드
  - 압축 파일(`.tar.gz`, Windows는 `.zip`) 생성
  - GitHub Release에 자동 업로드

### 릴리스 방법

1. 버전 태그 생성 및 푸시

```bash
git tag v0.2.0
git push origin v0.2.0
```

2. GitHub Actions에서 release job 완료 확인
3. GitHub Releases에서 OS별 아카이브 확인

## Release 후 터미널에서 바로 실행하기

릴리스 아카이브를 내려받아 `PATH`에 배치하면 어느 위치에서든 `gol-on-cli`를 실행할 수 있습니다.

### Linux / macOS

```bash
# 예시: Linux amd64 릴리스 설치
VERSION=v0.2.0
curl -fL -o gol-on-cli.tar.gz \
  "https://github.com/<OWNER>/<REPO>/releases/download/${VERSION}/gol-on-cli_${VERSION}_linux-amd64.tar.gz"

tar -xzf gol-on-cli.tar.gz
chmod +x gol-on-cli
sudo mv gol-on-cli /usr/local/bin/gol-on-cli

gol-on-cli --version
gol-on-cli --help
```

### Windows (PowerShell)

```powershell
$version = "v0.2.0"
Invoke-WebRequest -Uri "https://github.com/<OWNER>/<REPO>/releases/download/$version/gol-on-cli_${version}_windows-amd64.zip" -OutFile "gol-on-cli.zip"
Expand-Archive -Path .\gol-on-cli.zip -DestinationPath .\gol-on-cli

# PATH에 등록된 폴더(예: $env:USERPROFILE\bin)로 복사
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin" | Out-Null
Copy-Item .\gol-on-cli\gol-on-cli.exe "$env:USERPROFILE\bin\gol-on-cli.exe" -Force

# 필요 시 사용자 PATH에 $env:USERPROFILE\bin 추가 후 새 터미널 실행
& "$env:USERPROFILE\bin\gol-on-cli.exe" --version
```

> `<OWNER>/<REPO>`와 버전 값은 실제 릴리스에 맞게 바꿔서 사용하세요.
