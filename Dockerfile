FROM golang:1.24.6 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go test -v ./...
RUN go build -o app

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/app ./app
EXPOSE 8080
CMD ["./app"]
