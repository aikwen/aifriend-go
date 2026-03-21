# 编译
FROM golang:1.25 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/aifriend ./cmd/api

# 运行

FROM debian:bookworm-slim

WORKDIR /app

# 时区和证书
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Shanghai

COPY --from=builder /app/bin/aifriend /app/aifriend
COPY web /app/web

RUN mkdir -p /app/media/character/background_images /app/media/character/photos \
    && mkdir -p /app/media/user/photos

EXPOSE 8000
EXPOSE 8001

ENTRYPOINT ["/app/aifriend"]
CMD ["--config-file=/app/config/config.yaml"]