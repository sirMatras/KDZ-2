FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /file-storage ./cmd/file-storing-service

EXPOSE 8081

CMD ["/file-storage"]