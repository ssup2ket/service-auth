FROM golang:1.16 as builder
WORKDIR /workspace

# Build cnsenter
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ssup2ket-auth ./cmd/ssup2ket-auth

# Build image
FROM alpine:3.13.1
COPY --from=builder /workspace/ssup2ket-auth /usr/bin/ssup2ket-auth

CMD ["ssup2ket-auth"]
