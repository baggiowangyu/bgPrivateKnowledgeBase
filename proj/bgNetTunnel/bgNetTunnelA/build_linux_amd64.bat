SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build
rename bgNetTunnelA bgNetTunnelA_linux_amd64