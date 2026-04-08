FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY go.work ./
COPY shared ./shared
COPY services/main ./services/main
COPY services/auth ./services/auth
COPY services/room_manager ./services/room_manager
COPY services/room_worker ./services/room_worker
COPY services/vocab_manager ./services/vocab_manager

WORKDIR /workspace/services/room_manager
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/room_manager .

FROM alpine:3.23.3

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /out/room_manager ./room_manager

EXPOSE 8060

ENTRYPOINT ["./room_manager"]
