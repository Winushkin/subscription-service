FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs (или пропустить если нет)
RUN go install github.com/swaggo/swag/cmd/swag@latest || true
RUN swag init -g cmd/main.go || true

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]


