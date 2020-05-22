OUT_FILE=bin/kratos/pkg

install:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod download

test:
	go test -v ./...

run:
	ENV=dev go run main.go

build:
	go build -o ${OUT_FILE} -ldflags "-X main.gitCommit=$(shell git rev-list -1 HEAD)" main.go

