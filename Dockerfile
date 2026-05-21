FROM golang:1.26-alpine AS builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN swag init -g cmd/main.go

RUN CGO_ENABLED=0 go build -o app ./cmd


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/docs ./docs

EXPOSE 80

CMD ["./app"]
