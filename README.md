# ssup2ket-auth-service

## Development Environment

* Golang Version : 1.16

* Install protobuf-compiler

```
// Ubuntu
$ apt install -y protobuf-compiler

// MacOS
$ brew install protobuf
```

* Install oapi-codegen

```
$ go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.6.0
```

* Install gRPC binaries

```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

* Install mockery

```
$ go install github.com/vektra/mockery/v2@latest
```

* Install act

```
$ go install github.com/nektos/act@latest
```

## Reference

* UUID for DB - https://github.com/google/uuid/issues/20
* OpenTracing tracer - https://github.com/jaegertracing/jaeger-client-go/blob/master/zipkin/README.md
* OpenTracing middleware - https://github.com/go-chi/httptracer
