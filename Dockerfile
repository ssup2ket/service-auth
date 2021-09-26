# Build ssup2ket-auth binary
FROM golang:1.16 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ssup2ket-auth ./cmd/ssup2ket-auth

# Make up image
FROM alpine:3.13.1
WORKDIR /root
COPY --from=builder /workspace/ssup2ket-auth /usr/bin/ssup2ket-auth
COPY --from=builder /workspace/configs configs

CMD ["ssup2ket-auth"]
