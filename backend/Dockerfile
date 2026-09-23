FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/server .
COPY config.yaml .
EXPOSE 8080 9090
CMD ["./server"]
