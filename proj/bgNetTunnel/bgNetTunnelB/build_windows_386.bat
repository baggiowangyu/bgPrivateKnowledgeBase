SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build
rename bgNetTunnelB.exe bgNetTunnelB_windows_386.exe