FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mzt-api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

#RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

COPY --from=builder /app/mzt-api .

COPY ../.env .

ENTRYPOINT ["/app/mzt-api"]
