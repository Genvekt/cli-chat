LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.3.10

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../../.golangci.pipeline.yaml

generate-mocks:
	go generate ./internal/service/
	go generate ./internal/repository/
	go generate ./internal/config/

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/Genvekt/cli-chat/services/chat-server/internal/service/...,github.com/Genvekt/cli-chat/services/chat-server/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Genvekt/cli-chat/services/chat-server/internal/service/...,github.com/Genvekt/cli-chat/services/chat-server/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore