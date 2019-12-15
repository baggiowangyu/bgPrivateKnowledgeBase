SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build
rename bgNetTunnelA bgNetTunnelA_linux_386

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
rename bgNetTunnelA bgNetTunnelA_linux_amd64

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build
rename bgNetTunnelA.exe bgNetTunnelA_windows_386.exe

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build
rename bgNetTunnelA.exe bgNetTunnelA_windows_amd64.exe