clean:
	rm -rf ./bin

deploy: clean build-lambda
	serverless deploy --verbose --aws-profile serverless-deploy  

build-lambda:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/api cmd/lambda/main.go

build:
	go build -ldflags="-s -w" -o bin/api cmd/dev/main.go

test-unit:
	go test ./internal/...

test-unit-verbose:
	go test ./internal/... -v

test-behat:
	go test gherkin/api_test.go 

test-behat-verbose:
	go test -v gherkin/api_test.go 

run:
	air
