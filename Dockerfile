# Build service-auth binary
FROM golang:1.16 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o service-auth ./cmd/service-auth

# Make up image
FROM alpine:3.13.1
WORKDIR /root
COPY --from=builder /workspace/service-auth /usr/bin/service-auth
COPY --from=builder /workspace/configs configs

CMD ["service-auth"]
