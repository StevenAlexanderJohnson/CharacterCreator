FROM golang:1.25.0-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /character-creator ./cmd

FROM alpine:3.22

WORKDIR /

COPY --from=builder /character-creator .
COPY --from=builder /app/jwt_private_key.pem .
COPY --from=builder /app/dbmigration /dbmigration
COPY --from=builder /app/internal/templates /internal/templates
COPY --from=builder /app/public /public

EXPOSE 8080

CMD ["/character-creator"]