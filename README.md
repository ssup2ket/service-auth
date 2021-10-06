# ssup2ket-auth-service

* [Swagger](https://ssup2ket.github.io/ssup2ket-auth-service/api/openapi/swagger.html)
* [ER Diagram](https://drive.google.com/file/d/17gR4NP3bFl21aqhpr3PnhRePQTzafZoY/view?usp=sharing)

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

* Install yq

```
$ GO111MODULE=on go get github.com/mikefarah/yq/v4
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
* Istio GRPC - https://stackoverflow.com/questions/62459006/how-to-route-multiple-grpc-services-based-on-path-in-istiokubernetes
* Casbin RBAC - https://github.com/luk4z7/middleware-acl
* OpenTracing tracer - https://github.com/jaegertracing/jaeger-client-go/blob/master/zipkin/README.md
* OpenTracing middleware - https://github.com/go-chi/httptracer
* OpenTracing intercepter - https://github.com/opentracing-contrib/go-grpc
