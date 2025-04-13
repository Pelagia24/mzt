FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /app/.env .

COPY --from=builder /app/main .

# COPY --from=builder /app/cert.pem .
# COPY --from=builder /app/key.pem .

EXPOSE 443

CMD ["./main"]
