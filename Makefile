clean:
	go clean

deps:
	go get github.com/stretchr/testify/assert

test: clean
	go test -v

build:
	go clean
	go build
