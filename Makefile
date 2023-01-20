API_BINARY_NAME=api
RUNNER_BINARY_NAME=runner

build-api:
	GOOS=linux GOARCH=amd64 go build -o bin/${API_BINARY_NAME}-linux cmd/listener/listener.go
	GOOS=darwin GOARCH=arm64 go build -o bin/${API_BINARY_NAME}-darwin-apple-silicon cmd/listener/listener.go

build-runner:
	GOOS=linux GOARCH=amd64 go build -o bin/${RUNNER_BINARY_NAME}-linux cmd/runner/runner.go
	GOOS=darwin GOARCH=arm64 go build -o bin/${RUNNER_BINARY_NAME}-darwin-apple-silicon cmd/runner/runner.go

build: build-api build-runner

clean:
	go clean
	rm -f bin/${API_BINARY_NAME}-linux
	rm -f bin/${API_BINARY_NAME}-darwin-apple-silicon
	rm -f bin/${RUNNER_BINARY_NAME}-linux
	rm -f bin/${RUNNER_BINARY_NAME}-darwin-apple-silicon