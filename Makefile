clean:
	go clean

deps:
	go get github.com/stretchr/testify/assert

build:
	go clean
	go build
