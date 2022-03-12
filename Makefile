# Environment setup
env:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u github.com/golang/mock/mockgen

mocks:
	mockgen -source entities.go -destination migrationtest/provider.go -package migrationtest

test:
	go test ./... -v

test-integration:
	go test ./... -v -tags=integration

lint:
	golint -set_exit_status ./...