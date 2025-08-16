FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/http/main.go 

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/app ./app
COPY --from=builder /app/config/config.yaml ./config.yaml
COPY --from=builder /app/migrations/pgsql/* ./migrations/pgsql/

RUN apk add curl

CMD ["./app"]

EXPOSE 8000
