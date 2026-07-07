FROM golang:1.22-alpine AS builder
WORKDIR /src

COPY go.mod ./
COPY . .

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /out/server /app/server
COPY --from=builder /src/internal/db/migrations /app/internal/db/migrations

EXPOSE 8080 50051 8081

CMD ["/app/server"]
