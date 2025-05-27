FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /fileanalysisservice ./cmd/file-analysis-service

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /fileanalysisservice /fileanalysisservice

EXPOSE 8082

CMD ["/fileanalysisservice"]
