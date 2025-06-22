FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-gw ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api-gw .
COPY webadmin ./webadmin
COPY cert.pem .
COPY key.pem .
EXPOSE 8080
CMD ["./api-gw"]