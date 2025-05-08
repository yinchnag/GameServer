SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
SET GODEBUG=madvdontneed=1

go build -o app.exe -x -a ./zone/main.go