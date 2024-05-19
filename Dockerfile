FROM golang:alpine AS builder
RUN apk add --no-cache protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
WORKDIR /build
COPY . .
RUN protoc -I . proto/*.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:.
RUN go mod download
RUN go mod tidy
RUN go build -o /build/target .

FROM alpine:latest
COPY --from=builder /build/config/config-production.yaml ./config/config.yaml
COPY --from=builder /build/target .
CMD ["./target"]