# golang 1.16 version
.PHONY: all
all: test-unit run

# go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.6.0
# go install github.com/mikefarah/yq/v4@latest
.PHONY: gen-openapi
gen-openapi:
	oapi-codegen --generate types,chi-server,spec -o internal/server/http_server/http_server.gen.go --package http_server api/openapi/api.yml
	echo "var api_spec =" > api/openapi/api.json.js &&  yq eval -o=j api/openapi/api.yml >> api/openapi/api.json.js 

## Ubuntu 
# apt install protobuf-compiler
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
## MacOS
# brew install protobuf
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
.PHONY: gen-protobuf
gen-protobuf:
	protoc --go_out=. --go-grpc_out=. api/protobuf/api.proto

# go install github.com/vektra/mockery/v2@v2.9.0
.PHONY: gen-mock
gen-mock:
	mockery --all --dir internal/domain/repo --output internal/domain/repo/mocks
	mockery --all --dir internal/domain/service --output internal/domain/service/mocks

.PHONY: run
run: gen-openapi gen-protobuf
	. scripts/env-local && go run cmd/service-auth/main.go

.PHONY: init-local
init-local:
	docker run --name service-auth-local-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=local_auth -d mysql:8.0

.PHONY: build-image
build-image:
	docker build --tag ghcr.io/ssup2ket/service-auth:local .

# go install github.com/vektra/mockery/v2@v2.12.2
.PHONY: test-unit
test-unit: gen-mock
	go test -v -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o=cover.html

# Ubuntu
## apt install jq
## go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.6
# MacOS
## brew install jq
## go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.6
.PHONY: test-integration
test-integration:
	scripts/test-http-server.sh && scripts/test-grpc-server.sh

.PHONY: test-action
test-action: gen-mock
	act push --workflows ./.github/workflows/test-unit.yml -P ubuntu-20.04=ghcr.io/catthehacker/ubuntu:act-20.04
	act workflow_run --workflows ./.github/workflows/test-integration.yml -P ubuntu-20.04=ghcr.io/catthehacker/ubuntu:act-20.04
