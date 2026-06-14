$ErrorActionPreference = "Stop"

$CRONET_GO_VERSION = "dc1cda1fe28740ba069934ab62aeb8ef85388332"
$TAGS = "with_gvisor,with_quic,with_wireguard,with_utls,with_clash_api,with_grpc,with_awg,tfogo_checklinkname0,with_naive_outbound,with_conntrack"
$WINDOWS_ADD_TAGS = "with_purego"
$FINAL_TAGS = "$TAGS,$WINDOWS_ADD_TAGS"

$env:CGO_ENABLED = "1"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

if (Get-Command "x86_64-w64-mingw32-gcc" -ErrorAction SilentlyContinue) {
    $env:CC = "x86_64-w64-mingw32-gcc"
} else {
    $env:CC = "gcc"
}

Remove-Item -Recurse -Force bin\* -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path bin | Out-Null

go run -v "github.com/sagernet/cronet-go/cmd/build-naive@$CRONET_GO_VERSION" extract-lib --target windows/amd64 -o bin/

go build -trimpath -ldflags="-w -s -checklinkname=0 -buildid=" -buildmode=c-shared -tags $FINAL_TAGS -o bin/hiddify-core.dll ./platform/desktop

go install -mod=readonly github.com/akavel/rsrc@latest
go run ./cli tunnel exit -ErrorAction SilentlyContinue

Copy-Item bin/hiddify-core.dll ./hiddify-core.dll

$gopath = go env GOPATH
& "$gopath\bin\rsrc.exe" -ico ./assets/hiddify-cli.ico -o ./cmd/bydll/cli.syso

$env:CGO_LDFLAGS = "hiddify-core.dll"
go build -trimpath -ldflags="-w -s -checklinkname=0 -buildid=" -tags $TAGS -o bin/HiddifyCli.exe ./cmd/bydll

Remove-Item ./hiddify-core.dll -ErrorAction SilentlyContinue

echo "Successfully built hiddify-core for windows."
