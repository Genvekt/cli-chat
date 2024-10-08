FROM golang:1.22.5-alpine as builder

COPY libraries /app/libraries

COPY services/chat_server /app/services/chat_server

WORKDIR /app/services/chat_server

RUN go mod download
RUN go build -o ./bin/main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/services/chat_server/bin/main ./main
COPY --from=builder /app/services/chat_server/local.env ./.env

CMD ["./main"]