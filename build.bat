@REM 执行便可打包为armv7


 SET CGO_ENABLED=0

 SET GOOS=linux

 SET GOARCH=arm

 SET GOARM=7

go build -o ./app_linux_armv7 main.go

@REM 打包为amd64
@REM set GOARCH=amd64
@REM set GOOS=linux
@REM go build main.go


