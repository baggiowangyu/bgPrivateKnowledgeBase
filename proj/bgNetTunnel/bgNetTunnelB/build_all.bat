echo "cleanup files"
DEL bgNetTunnelB
DEL bgNetTunnelB.exe
DEL bgNetTunnelB_linux_386
DEL bgNetTunnelB_linux_amd64
DEL bgNetTunnelB_windows_386.exe
DEL bgNetTunnelB_windows_amd64.exe

mkdir bin
cd bin

set hour=%time:~0,2%
if /i %hour% LSS 10 (
 set hour=0%time:~1,1%
)
set filename=%date:~0,4%%date:~5,2%%date:~8,2%_%hour%.%time:~3,2%.%time:~6,2%
echo "´´½¨Ä¿Â¼£º%filename%"
mkdir %filename%
cd %filename%
mkdir config

cd ..
cd ..

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build
rename bgNetTunnelB bgNetTunnelB_linux_386

COPY "bgNetTunnelB_linux_386" "bin/%filename%/bgNetTunnelB_linux_386" /Y

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
rename bgNetTunnelB bgNetTunnelB_linux_amd64

COPY "bgNetTunnelB_linux_amd64" "bin/%filename%/bgNetTunnelB_linux_amd64" /Y

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_386.exe

COPY "bgNetTunnelB_windows_386.exe" "bin/%filename%/bgNetTunnelB_windows_386.exe" /Y

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_amd64.exe

COPY "bgNetTunnelB_windows_amd64.exe" "bin/%filename%/bgNetTunnelB_windows_amd64.exe" /Y

cd config

COPY "config.toml" "../bin/%filename%/config/config.toml" /Y

cd ..