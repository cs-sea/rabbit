linux:
	@go mod tidy
	@go mod vendor
	@GOOS=linux GOARCH=amd64 go build -ldflags "-w -s"  -v -o rabbit main/main.go

mac:
	@go mod tidy
	@go mod vendor
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s" -v -o rabbit main/main.go

test:
	go test -v ./...
