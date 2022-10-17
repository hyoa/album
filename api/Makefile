clean:
	rm -rf ./bin

deploy: clean build
	serverless deploy --verbose

build:
	go build -ldflags="-s -w" -o bin/api cmd/dev/main.go

test-unit:
	go test ./intenal/...

test-unit-verbose:
	go test ./intenal/... -v

test-behat:
	go test gherkin/api_test.go 

test-behat-verbose:
	go test -v gherkin/api_test.go 

run:
	air