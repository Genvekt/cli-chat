LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.3.10
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../../.golangci.pipeline.yaml

docker-build:
	docker build --no-cache --platform linux/amd64 -t auth -f ./local.Dockerfile ../../

generate-mocks:
	go generate ./internal/service/
	go generate ./internal/repository/
	go generate ./internal/config/
	go generate ./internal/utils/

generate-swagger:
	cp ../../libraries/api/user/swagger/api.swagger.json ./pkg/swagger/userApi.swagger.json
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/Genvekt/cli-chat/services/auth/internal/service/...,github.com/Genvekt/cli-chat/services/auth/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/Genvekt/cli-chat/services/auth/internal/service/...,github.com/Genvekt/cli-chat/services/auth/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

gen-cert:
	openssl genrsa -out certificate/keys/ca.key 4096
	openssl req -new -x509 -key certificate/keys/ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out certificate/keys/ca.cert
	openssl genrsa -out certificate/keys/service.key 4096
	openssl req -new -key certificate/keys/service.key -out certificate/keys/service.csr -config certificate/certificate.conf
	openssl x509 -req -in certificate/keys/service.csr -CA certificate/keys/ca.cert -CAkey certificate/keys/ca.key -CAcreateserial \
    		-out certificate/keys/service.pem -days 365 -sha256 -extfile certificate/certificate.conf -extensions req_ext