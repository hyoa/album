clean:
	rm -rf ./bin

deploy: clean build
	serverless deploy --verbose

build:
	go build -ldflags="-s -w" -o bin/api cmd/dev/main.go

test:
	go test ./...

test-verbose:
	go test ./... -v

run:
	air