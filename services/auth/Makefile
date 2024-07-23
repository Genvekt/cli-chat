LOCAL_BIN:=$(CURDIR)/bin
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../../.golangci.pipeline.yaml

docker-build:
	docker build --no-cache --platform linux/amd64 -t auth -f ./local.Dockerfile ../../