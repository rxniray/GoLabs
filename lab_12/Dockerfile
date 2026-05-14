FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./${SERVICE_NAME}

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
