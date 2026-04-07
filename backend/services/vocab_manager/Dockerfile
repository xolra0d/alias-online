FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY go.work ./
COPY shared ./shared
COPY services/main ./services/main
COPY services/auth ./services/auth
COPY services/vocabs ./services/vocabs

WORKDIR /workspace/services/vocabs
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/vocab_manager .

FROM alpine:3.23.3

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /out/vocabs ./vocabs

EXPOSE 8079

ENTRYPOINT ["./vocabs"]
