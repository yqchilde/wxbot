FROM golang:1.19-alpine as builder

ENV GOPROXY="https://goproxy.cn,direct"
ARG VERSION
WORKDIR /app/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/yqchilde/wxbot/engine/robot.version=${VERSION}" -o ./bin/wxbot

FROM alpine:latest
LABEL MAINTAINER="yqchilde@gmail.com"
WORKDIR /app/
VOLUME /app/data/

RUN apk add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/shanghai" >> /etc/timezone \
    && apk del tzdata

COPY --from=builder /app/bin/wxbot /app/wxbot
COPY --from=builder /app/config.yaml /app/config.yaml

EXPOSE 9528
ENTRYPOINT ["./wxbot"]
