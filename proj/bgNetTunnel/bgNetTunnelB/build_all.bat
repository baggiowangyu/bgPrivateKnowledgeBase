SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build
rename bgNetTunnelB bgNetTunnelB_linux_386

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
rename bgNetTunnelB bgNetTunnelB_linux_amd64

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_386.exe

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_amd64.exe