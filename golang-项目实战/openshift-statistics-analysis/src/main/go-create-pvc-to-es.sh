::gopath 为项目目录
SET GO15VENDOREXPERIMENT=1
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o go-create-pvc-to-es go-create-pvc-to-es.go