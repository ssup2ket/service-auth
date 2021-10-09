# ssup2ket-auth-service

ssup2ket-auth-service is the service responsible for user management and authentication/authorization in the [ssup2ket](https://github.com/ssup2ket/ssup2ket) Project. ssup2ket-auth-service follows this [considerations](https://github.com/ssup2ket/ssup2ket#ssup2ket-service-considerations).

* [Swagger](https://ssup2ket.github.io/ssup2ket-auth-service/api/openapi/swagger.html)
* [ER Diagram](https://drive.google.com/file/d/17gR4NP3bFl21aqhpr3PnhRePQTzafZoY/view?usp=sharing)

## Authentication/Authorization

ssup2ket-auth-service uses simple authentication based on **ID/Password**. A user can get the **JWT** Token required for authentication/authorization by entering ID/Password. Passwords are encrypted and stored using the PBKDF2 algorithm.

In JWT Token, **User's ID(UUID), Login ID, Password and Role** are stored. Other services of the ssup2ket Project need to implement authentication and RBAC-based authorization through JWT Token. Each User can have only one Role. There are two roles, admin and user.

## Used main external packages and tools

ssup2ket-auth-service uses following external packages and tools.

* **HTTP Server, Middleware** - [chi](https://github.com/go-chi/chi), [HTTP](https://pkg.go.dev/net/http), [oapi-codegen](https://github.com/deepmap/oapi-codegen)
* **GRPC Server, Intercepter** - [grpc](https://pkg.go.dev/google.golang.org/grpc), [protoc-gen-go](https://pkg.go.dev/github.com/golang/protobuf/protoc-gen-go)
* **MySQL** - [GORM](https://gorm.io/index.html)
* **Kafaka** - [kafka-go](https://github.com/segmentio/kafka-go)
* **Authorziation** - [Casbin](https://casbin.org/)
* **Logging, Tracking** - [zerolog](https://github.com/rs/zerolog), [Istio](https://istio.io/)
* **CI/CD** - [Testify](https://github.com/stretchr/testify), [sqlmock](https://github.com/DATA-DOG/go-sqlmock), [Mockery](https://github.com/mockery/mockery), [Github Actions](https://github.com/features/actions), [K8s](https://kubernetes.io/), [ArgoCD](https://argo-cd.readthedocs.io/en/stable/), [ArgoCD Image Updater](https://github.com/argoproj-labs/argocd-image-updater), [Kustomize](https://kustomize.io/)

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
