SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_amd64.exe