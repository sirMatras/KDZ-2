FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . ./

RUN go build -o /api-gateway ./cmd/gateway

EXPOSE 8080

CMD ["/api-gateway"]
