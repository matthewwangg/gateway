FROM golang:1.23.0-alpine

WORKDIR /gateway

COPY go.mod ./

RUN  go mod download

COPY . .

RUN go build -o gateway-server ./cmd/main.go

CMD ["./gateway-server"]
