set GOPATH="%cd%"
mkdir bin
go build -o bin\lightproxy.exe src\main\startup.go
copy src\main\config.json bin\