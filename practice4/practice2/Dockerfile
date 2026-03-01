FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o myapp ./cmd/api/main.go

# Run
FROM alpine
WORKDIR /app
COPY --from=builder /app/myapp .
COPY --from=builder /app/database ./database
ENTRYPOINT [ "./myapp" ]
