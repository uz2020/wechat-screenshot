all: main.go
	GOOS=windows GOARCH=amd64 go build .
