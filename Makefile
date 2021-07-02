# golang 1.16 version
.PHONY: all
all: test run

# go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.6.0
.PHONY: gen-openapi
gen-openapi:
	oapi-codegen --generate types,chi-server,spec -o internal/server/http_server/http_server.gen.go --package http_server api/openapi/api.yml

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
	source scripts/local-run-env.sh && go run cmd/ssup2ket-auth/main.go

.PHONY: test
test: gen-mock
	go test -v ./...