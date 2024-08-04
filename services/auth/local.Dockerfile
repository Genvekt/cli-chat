FROM golang:1.22.5-alpine as builder

COPY libraries /app/libraries

COPY services/auth /app/services/auth

WORKDIR /app/services/auth

RUN go mod download
RUN go build -o ./bin/main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/services/auth/bin/main ./main
COPY --from=builder /app/services/auth/local.env ./.env

CMD ["./main"]