linux:
	@sh shell/linux.sh

mac:
	@sh shell/mac.sh

test:
	go test -v ./...
