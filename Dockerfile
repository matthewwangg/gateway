FROM golang:1.23.0 AS builder

WORKDIR /gateway

RUN apt-get update && apt-get install -y --no-install-recommends \
    libprotobuf-dev protobuf-compiler \
    libgrpc++-dev protobuf-compiler-grpc \
    && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod ./

RUN go mod download

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/pb/log.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-server ./cmd/main.go

FROM alpine:latest

WORKDIR /gateway

COPY --from=builder /gateway/gateway-server /gateway/gateway-server

CMD ["./gateway-server"]
