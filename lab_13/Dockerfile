FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]